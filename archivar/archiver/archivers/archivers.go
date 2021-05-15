package archivers

import (
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type factory func(c interface{}, logger *logrus.Logger) Archiver

var registered = make(map[string]factory)

// UploadFunc takes a file and uploads it to their archive backend
type UploadFunc func(file.File) (err error)

// Archiver is used to store files in their archive backend
type Archiver interface {
	Upload(file.File) (err error)
	Connect() (err error)
}

// Register a new Archiver
func Register(p factory) {
	registered[caller.FactoryPackage()] = p
}

// Get a registered archiver
func Get(n string, c interface{}, logger *logrus.Logger) Archiver {
	p, exists := registered[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}
