package filters

import (
	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, logger *logrus.Logger) Filter

var registeredFilters = make(map[string]factory)

// Filter will return filterResult.Results and cause rejected to be not further processed
type Filter interface {
	Filter(*file.File) (filterResult.Results, error)
}

// Register a new filter
func Register(p factory) {
	registeredFilters[caller.FactoryPackage()] = p
}

// Get a registered filter
func Get(n string, c interface{}, logger *logrus.Logger) Filter {
	p, exists := registeredFilters[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}

func ListFilters() (n []string) {
	for f := range registeredFilters {
		n = append(n, f)
	}

	return
}
