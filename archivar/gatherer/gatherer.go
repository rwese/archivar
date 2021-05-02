package gatherer

import (
	"errors"

	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/gatherer/imap"
	"github.com/sirupsen/logrus"
)

type Gatherer interface {
	Connect() (err error)
	Download() (err error)
}

type GathererConfig struct {
	Type         string
	Server       string
	Username     string
	Password     string
	Token        string
	KeepUploaded bool
	ClientId     string
	ClientSecret string
}

func New(g GathererConfig, archivar archiver.Archiver, logger *logrus.Logger) (gatherer Gatherer, err error) {
	switch g.Type {
	case "imap":
		gatherer = imap.New(
			g.Server,
			g.Username,
			g.Password,
			g.KeepUploaded,
			archivar,
			logger,
		)
	default:
		err = errors.New("could not create new gatherer from given config")
	}

	return
}
