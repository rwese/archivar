package gatherer

import (
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/gatherer/gatherers/imap"
	"github.com/sirupsen/logrus"
)

type Gatherer interface {
	Connect() error
	Download() error
}

func New(gathererType string, config interface{}, archivar archiver.Archiver, logger *logrus.Logger) Gatherer {
	switch gathererType {
	case "imap":
		return imap.New(config, archivar, logger)
	default:
		logger.Panic("could not create new gatherer from given config")
	}

	return nil
}
