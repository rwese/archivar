package archivers

import (
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, logger *logrus.Logger) Archiver

var registered = make(map[string]factory)

type UploadFunc func(file.File) (err error)

type Archiver interface {
	Upload(file.File) (err error)
	Connect() (err error)
}

func Register(p factory) {
	registered[caller.FactoryPackage()] = p
}

func Get(n string, c interface{}, logger *logrus.Logger) Archiver {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}
