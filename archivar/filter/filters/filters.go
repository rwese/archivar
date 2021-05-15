package filters

import (
	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, logger *logrus.Logger) Filter

var registered = make(map[string]factory)

type UploadFunc func(file.File) (err error)

type Filter interface {
	Filter(*file.File) (filterResult.Results, error)
}

func Register(p factory) {
	registered[caller.FactoryPackage()] = p
}

func Get(n string, c interface{}, logger *logrus.Logger) Filter {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}
