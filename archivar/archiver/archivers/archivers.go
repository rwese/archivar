package archivers

import (
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/caller"
	"github.com/sirupsen/logrus"
)

type archiverFactory func(c interface{}, logger *logrus.Logger) Archiver

var registeredArchiver = make(map[string]archiverFactory)

// UploadFunc takes a file and uploads it to their archive backend
type UploadFunc func(*file.File) (err error)

// Archiver is used to store files in their archive backend
type Archiver interface {
	Upload(*file.File) (err error)
	Connect() (err error)
}

// Register a new Archiver
func RegisterArchiver(p archiverFactory) {
	registeredArchiver[caller.FactoryPackage()] = p
}

// Get a registered archiver
func GetArchiver(n string, c interface{}, logger *logrus.Logger) Archiver {
	p, exists := registeredArchiver[n]
	if !exists {
		return nil
	}

	return p(c, logger)
}

func ListArchivers() (archiverNames []string) {
	for a := range registeredArchiver {
		archiverNames = append(archiverNames, a)
	}

	return
}
