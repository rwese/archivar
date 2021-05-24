package brains

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type brainFactory func(c interface{}, archiver archivers.Archiver, gatherer archivers.Gatherer, logger *logrus.Logger) Brain

var registeredBrain = make(map[string]brainFactory)

// Brain is used to download files and give them to their storage
type Brain interface {
	Connect() error
	DoYouKnow(stuff ...interface{}) bool
	Remember(stuff ...interface{})
}

// Register a new brain
func RegisterBrain(p brainFactory) {
	registeredBrain[caller.FactoryPackage()] = p
}

// Get a brain from the registry
func GetBrain(n string, c interface{}, archiver archivers.Archiver, gatherer archivers.Gatherer, logger *logrus.Logger) Brain {

	p, exists := registeredBrain[n]
	if !exists {
		return nil
	}

	return p(c, archiver, gatherer, logger)
}

func ListBrains() (n []string) {
	for g := range registeredBrain {
		n = append(n, g)
	}

	return
}
