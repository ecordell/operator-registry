package mirror

import (
	"context"
	"database/sql"
	"strings"

	"github.com/docker/distribution/reference"
	"k8s.io/apimachinery/pkg/util/errors"

	"github.com/operator-framework/operator-registry/pkg/sqlite"
)

type Mirrorer interface {
	Mirror() (map[string]string, error)
}

// DatabaseExtractor knows how to pull an index image and extract its database
type DatabaseExtractor interface {
	Extract(from string) (string, error)
}

type DatabaseExtractorFunc func(from string) (string, error)

func (f DatabaseExtractorFunc) Extract(from string) (string, error) {
	return f(from)
}

// ImageMirrorer knows how to mirror an image from one registry to another
type ImageMirrorer interface {
	Mirror(from, to string) error
}

type ImageMirrorerFunc func(from, to string) error

func (f ImageMirrorerFunc) Mirror(from, to string) error {
	return f(from, to)
}

type IndexImageMirrorer struct {
	ImageMirrorer       ImageMirrorer
	DatabaseExtractor   DatabaseExtractor

	// options
	Source, Dest        string
}

var _ Mirrorer = &IndexImageMirrorer{}

func NewIndexImageMirror(config *IndexImageMirrorerOptions, options ...ImageIndexMirrorOption) (*IndexImageMirrorer, error) {
	if config == nil {
		config = DefaultImageIndexMirrorerOptions()
	}
	config.Apply(options)
	if err := config.Complete(); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &IndexImageMirrorer{
		ImageMirrorer:         config.ImageMirrorer,
		DatabaseExtractor:     config.DatabaseExtractor,
		Source:                config.Source,
		Dest:                  config.Dest,
	}, nil
}

func (b *IndexImageMirrorer) Mirror() (map[string]string, error) {
	dbPath, err := b.DatabaseExtractor.Extract(b.Source)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	migrator, err := sqlite.NewSQLLiteMigrator(db)
	if err != nil {
		return nil, err
	}
	if err := migrator.Migrate(context.TODO()); err != nil {
		return nil, err
	}

	querier := sqlite.NewSQLLiteQuerierFromDb(db)
	images, err := querier.ListImages(context.TODO())
	if err != nil {
		return nil, err
	}

	mapping := map[string]string{}

	errs := []error{}

	for _, img := range images {
		ref, err := reference.ParseNamed(img)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		i := strings.IndexRune(ref.Name(), '/')
		if i == -1 || (!strings.ContainsAny(ref.Name()[:i], ":.") && ref.Name()[:i] != "localhost") {
			continue
		}
		name := ref.Name()[i+1:]

		mapping[img] = strings.Join([]string{b.Dest, name}, "/")
	}

	for from, to := range mapping {
		if err := b.ImageMirrorer.Mirror(from, to); err != nil {
			errs = append(errs, err)
		}
	}

	return mapping, errors.NewAggregate(errs)
}
