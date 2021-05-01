package archiver

import (
	"errors"
	"io"

	"github.com/rwese/archivar/archivar/archiver/webdav"
	"github.com/sirupsen/logrus"
)

type ArchiverConfig struct {
	Type            string
	Server          string
	Username        string
	Password        string
	UploadDirectory string
}

type Archiver interface {
	Upload(fileName string, directory string, fileHandle io.Reader) (err error)
}

func New(a ArchiverConfig, logger *logrus.Logger) (archiver Archiver, err error) {

	switch a.Type {
	case "webdav":
		archiver = webdav.New(
			a.Server,
			a.Username,
			a.Password,
			a.UploadDirectory,
			logger,
		)
	default:
		err = errors.New("could not create new archiver from given config")
	}

	return
}
