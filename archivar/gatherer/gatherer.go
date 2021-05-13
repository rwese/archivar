package gatherer

import (
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/gatherer/gatherers/imap"
	"github.com/rwese/archivar/archivar/gatherer/gatherers/webdav"
	"github.com/sirupsen/logrus"
)

type Gatherer interface {
	Download() error
}

func New(gathererType string, config interface{}, archivar archiver.Archiver, logger *logrus.Logger) Gatherer {
	switch gathererType {
	case "imap":
		return imap.New(config, archivar, logger)
	case "webdav":
		return webdav.New(config, archivar, logger)
	default:
		logger.Panicf("could not create new gatherer '%s' from given config", gathererType)
	}

	return nil
}
