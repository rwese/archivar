package gatherers

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, storage archivers.Archiver, logger *logrus.Logger) Gatherer

var registered = make(map[string]factory)

// Gatherer is used to download files and give them to their storage
type Gatherer interface {
	Download() error
	Connect() (err error)
}

// Register a new gatherer
func Register(p factory) {
	registered[caller.FactoryPackage()] = p
}

// Get a gatherer from the registry
func Get(n string, c interface{}, storage archivers.Archiver, logger *logrus.Logger) Gatherer {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, storage, logger)
}
