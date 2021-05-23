package archivers

import (
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type gathererFactory func(c interface{}, storage Archiver, logger *logrus.Logger) Gatherer

var registeredGatherer = make(map[string]gathererFactory)

// Gatherer is used to download files and give them to their storage
type Gatherer interface {
	Download() error
	Connect() (err error)
}

// Register a new gatherer
func RegisterGatherer(p gathererFactory) {
	registeredGatherer[caller.FactoryPackage()] = p
}

// Get a gatherer from the registry
func GetGatherer(n string, c interface{}, storage Archiver, logger *logrus.Logger) Gatherer {
	p, exists := registeredGatherer[n]
	if !exists {
		return nil
	}

	return p(c, storage, logger)
}
