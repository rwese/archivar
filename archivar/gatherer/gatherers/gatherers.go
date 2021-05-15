package gatherers

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/utils/caller"
	"github.com/sirupsen/logrus"
)

type Factory func(c interface{}, storage archivers.Archiver, logger *logrus.Logger) Gatherer

var registered = make(map[string]Factory)

type Gatherer interface {
	Download() error
}

func Register(p Factory) {
	registered[caller.FactoryPackage()] = p
}

func Get(n string, c interface{}, storage archivers.Archiver, logger *logrus.Logger) Gatherer {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, storage, logger)
}
