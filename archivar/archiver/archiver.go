package archiver

import (
	"github.com/rwese/archivar/archivar/archiver/archivers/google_drive"
	"github.com/rwese/archivar/archivar/archiver/archivers/webdav"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

type UploadFunc func(file.File) (err error)
type Archiver interface {
	Upload(file.File) (err error)
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
