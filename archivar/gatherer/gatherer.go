package gatherer

import (
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/filter"
	"github.com/rwese/archivar/archivar/gatherer/imap"
	"github.com/sirupsen/logrus"
)

type Gatherer interface {
	Connect() (err error)
	Download() (err error)
}

func New(gathererType string, config interface{}, archivar archiver.Archiver, filters []filter.Filter, logger *logrus.Logger) Gatherer {
	switch gathererType {
	case "imap":
		return imap.New(config, archivar, filters, logger)
	default:
		logger.Panic("could not create new gatherer from given config")
	}

	return nil
}
