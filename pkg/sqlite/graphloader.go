package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/operator-framework/operator-registry/pkg/registry"
)

// GraphLoader generates a graph
// GraphLoader supports multiple different loading schemes
// GraphLoader from SQL, GraphLoader from old format (filesystem), GraphLoader from SQL + input bundles
type GraphLoader interface {
	Generate() (*registry.Package, error)
}

type SQLGraphLoader struct {
	Querier     *SQLQuerier
	PackageName string
}

type ChannelEntryNode struct {
	PackageName        string
	ChannelName        string
	BundleName         string
	BundlePath         string
	Version            string
	Replaces           string
	ReplacesVersion    string
	ReplacesBundlePath string
}

func NewSQLGraphLoader(dbFilename, name string) (*SQLGraphLoader, error) {
	querier, err := NewSQLLiteQuerier(dbFilename)
	if err != nil {
		return nil, err
	}

	return &SQLGraphLoader{
		Querier:     querier,
		PackageName: name,
	}, nil
}

func NewSQLGraphLoaderFromDB(db *sql.DB, name string) (*SQLGraphLoader, error) {
	return &SQLGraphLoader{
		Querier:     NewSQLLiteQuerierFromDb(db),
		PackageName: name,
	}, nil
}

func (g *SQLGraphLoader) Generate() (*registry.Package, error) {
	ctx := context.TODO()
	defaultChannel, err := g.Querier.GetDefaultPackage(ctx, g.PackageName)
	if err != nil {
		return nil, err
	}

	channelEntries, err := g.Querier.GetChannelEntriesFromPackage(ctx, g.PackageName)
	if err != nil {
		return nil, err
	}

	channels, err := g.GraphFromEntries(channelEntries)
	if err != nil {
		return nil, err
	}

	return &registry.Package{
		Name:           g.PackageName,
		DefaultChannel: defaultChannel,
		Channels:       channels,
	}, nil
}

// GraphFromEntries builds the graph from a set of channel entries
func (g *SQLGraphLoader) GraphFromEntries(channelEntries []ChannelEntryNode) (map[string]registry.Channel, error) {
	channels := map[string]registry.Channel{}

	type replaces map[registry.BundleKey]map[registry.BundleKey]struct{}

	channelGraph := map[string]replaces{}
	channelHeadCandidates := map[string]map[registry.BundleKey]struct{}{}

	// add all channels and nodes to the graph
	for _, entry := range channelEntries {
		// create channel if we haven't seen it yet
		if _, ok := channelGraph[entry.ChannelName]; !ok {
			channelGraph[entry.ChannelName] = replaces{}
		}

		key := registry.BundleKey{
			BundlePath: entry.BundlePath,
			Version:    entry.Version,
			CsvName:    entry.BundleName,
		}
		channelGraph[entry.ChannelName][key] = map[registry.BundleKey]struct{}{}

		// every bundle in a channel is a potential head of that channel
		if _, ok := channelHeadCandidates[entry.ChannelName]; !ok {
			channelHeadCandidates[entry.ChannelName] = map[registry.BundleKey]struct{}{key: {}}
		} else {
			channelHeadCandidates[entry.ChannelName][key] = struct{}{}
		}
	}

	// add all edges to the graph
	for _, entry := range channelEntries {
		key := registry.BundleKey{
			BundlePath: entry.BundlePath,
			Version:    entry.Version,
			CsvName:    entry.BundleName,
		}
		replacesKey := registry.BundleKey{
			BundlePath: entry.BundlePath,
			Version:    entry.ReplacesVersion,
			CsvName:    entry.Replaces,
		}

		if !replacesKey.IsEmpty() {
			channelGraph[entry.ChannelName][key][replacesKey] = struct{}{}
		}

		delete(channelHeadCandidates[entry.ChannelName], replacesKey)
	}

	for channelName, candidates := range channelHeadCandidates {
		if len(candidates) == 0 {
			return nil, fmt.Errorf("no channel head found for %s", channelName)
		}
		if len(candidates) > 1 {
			return nil, fmt.Errorf("multiple candidate channel heads found for %s: %v", channelName, candidates)
		}

		for head := range candidates {
			channel := registry.Channel{
				Head:     head,
				Replaces: channelGraph[channelName],
			}
			channels[channelName] = channel
		}
	}

	return channels, nil
}
