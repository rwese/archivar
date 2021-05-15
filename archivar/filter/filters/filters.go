package filters

import (
	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/utils/caller"
	"github.com/sirupsen/logrus"
)

type Factory func(c interface{}, logger *logrus.Logger) Filter

var registered = make(map[string]Factory)

type UploadFunc func(file.File) (err error)

type Filter interface {
	Filter(*file.File) (filterResult.Results, error)
}

func Register(p Factory) {
	registered[caller.FactoryPackage()] = p
}

func Get(n string, c interface{}, logger *logrus.Logger) Filter {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}
