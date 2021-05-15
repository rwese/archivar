package processors

import (
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/utils/caller"
	"github.com/sirupsen/logrus"
)

type Factory func(c interface{}, logger *logrus.Logger) Processor

var registered = make(map[string]Factory)

type Processor interface {
	Process(*file.File) error
}

func Register(p Factory) {
	registered[caller.FactoryPackage()] = p
}

func Get(n string, c interface{}, logger *logrus.Logger) Processor {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}
