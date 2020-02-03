package registry

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// PackageRemover removes a package from the database
type PackageRemover struct {
	loader   Load
	packages string
}

var _ Remover = &PackageRemover{}

func NewRemoverForPackages(loader Load, packages string) *PackageRemover {
	return &PackageRemover{
		loader:   loader,
		packages: packages,
	}
}

func (d *PackageRemover) Remove() error {
	log := logrus.WithField("pkg", d.packages)

	log.Info("deleting packages")

	var errs []error
	packages := sanitizePackageList(strings.Split(d.packages, ","))
	log.Info("input has been sanitized")
	log.Infof("packages: %s", packages)

	for _, pkg := range packages {
		if err := d.loader.RmPackageName(pkg); err != nil {
			errs = append(errs, fmt.Errorf("error removing operator package %s: %s", pkg, err))
		}
	}

	return utilerrors.NewAggregate(errs)
}

// sanitizePackageList sanitizes the set of package(s) specified. It removes
// duplicates and ignores empty string.
func sanitizePackageList(in []string) []string {
	out := make([]string, 0)

	inMap := map[string]bool{}
	for _, item := range in {
		if _, ok := inMap[item]; ok || item == "" {
			continue
		}

		inMap[item] = true
		out = append(out, item)
	}

	return out
}
