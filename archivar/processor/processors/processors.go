package processors

import (
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, logger *logrus.Logger) Processor

var registered = make(map[string]factory)

type Processor interface {
	Process(*file.File) error
}

func Register(p factory) {
	registered[caller.FactoryPackage()] = p
}

func Get(n string, c interface{}, logger *logrus.Logger) Processor {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}
