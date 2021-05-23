package processors

import (
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, logger *logrus.Logger) Processor

var registeredProcessors = make(map[string]factory)

// Processor will modify the given file
type Processor interface {
	Process(*file.File) error
}

// Register a new processor
func Register(p factory) {
	registeredProcessors[caller.FactoryPackage()] = p
}

// Get a processor from the registry
func Get(n string, c interface{}, logger *logrus.Logger) Processor {
	p, exists := registeredProcessors[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}

func ListProcessors() (n []string) {
	for f := range registeredProcessors {
		n = append(n, f)
	}

	return
}
