package archiver

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	_ "github.com/rwese/archivar/archivar/archiver/archivers/google_drive"
	_ "github.com/rwese/archivar/archivar/archiver/archivers/webdav"
	"github.com/sirupsen/logrus"
)

func New(typeName string, config interface{}, logger *logrus.Logger) (archiver archivers.Archiver) {
	g := archivers.Get(typeName, config, logger)

	if g == nil {
		logger.Panicf("could not create new archiver '%s' from given config", typeName)
	}

	return g
}
