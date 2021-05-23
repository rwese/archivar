package gatherer

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/gatherer/gatherers"
	"github.com/sirupsen/logrus"

	_ "github.com/rwese/archivar/archivar/gatherer/gatherers/filesystem"
	_ "github.com/rwese/archivar/archivar/gatherer/gatherers/imap"
	_ "github.com/rwese/archivar/archivar/gatherer/gatherers/webdav"
)

// New will return a new gatherer based on the given typeName and config
func New(typeName string, config interface{}, archivar archivers.Archiver, logger *logrus.Logger) gatherers.Gatherer {
	g := gatherers.Get(typeName, config, archivar, logger)

	if g == nil {
		logger.Panicf("could not create new gatherer '%s' from given config", typeName)
	}

	return g
}
