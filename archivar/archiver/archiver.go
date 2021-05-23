package archiver

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	_ "github.com/rwese/archivar/archivar/archiver/archivers/filesystem"
	_ "github.com/rwese/archivar/archivar/archiver/archivers/google_drive"
	_ "github.com/rwese/archivar/archivar/archiver/archivers/webdav"
	"github.com/sirupsen/logrus"
)

// NewArchiver will return a new archiver backend based on the given typeName and config
func NewArchiver(typeName string, config interface{}, logger *logrus.Logger) archivers.Archiver {
	archiver := archivers.GetArchiver(typeName, config, logger)

	if archiver == nil {
		logger.Panicf("could not create new archiver '%s' from given config", typeName)
	}

	return archiver
}

// NewArchiver will return a new archiver backend based on the given typeName and config
func NewGatherer(typeName string, config interface{}, storage archivers.Archiver, logger *logrus.Logger) archivers.Gatherer {
	gatherer := archivers.GetGatherer(typeName, config, storage, logger)

	if gatherer == nil {
		logger.Panicf("could not create new gatherer '%s' from given config", typeName)
	}

	return gatherer
}
