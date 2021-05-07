package archiver

import (
	"errors"
	"io"

	"github.com/rwese/archivar/archivar/archiver/google_drive"
	"github.com/rwese/archivar/archivar/archiver/webdav"
	"github.com/sirupsen/logrus"
)

type ArchiverConfig struct {
	Type            string
	Server          string
	Username        string
	Password        string
	UploadDirectory string
	OAuthToken      string
	ClientSecrets   string
}

type Archiver interface {
	Connect() (newSession bool, err error)
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
	case "gdrive":
		archiver = google_drive.New(
			a.OAuthToken,
			a.ClientSecrets,
			a.UploadDirectory,
			logger,
		)
	default:
		err = errors.New("could not create new archiver from given config")
	}

	return
}
