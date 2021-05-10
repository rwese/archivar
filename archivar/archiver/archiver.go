package archiver

import (
	"io"

	"github.com/rwese/archivar/archivar/archiver/archivers/google_drive"
	"github.com/rwese/archivar/archivar/archiver/archivers/webdav"
	"github.com/sirupsen/logrus"
)

type Archiver interface {
	Connect() (newSession bool, err error)
	Upload(fileName string, directory string, fileHandle io.Reader) (err error)
}

func New(archiverType string, config interface{}, logger *logrus.Logger) (archiver Archiver) {
	switch archiverType {
	case "webdav":
		return webdav.New(config, logger)
	case "gdrive":
		return google_drive.New(config, logger)
	default:
		logger.Panicf("could not create new archiver '%s' from given config", archiverType)
	}

	return nil
}
