package webdav

import (
	"path"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/internal/file"
	webdavClient "github.com/rwese/archivar/internal/webdav"
	"github.com/rwese/archivar/utils/config"

	"github.com/sirupsen/logrus"
)

// Webdav allows to upload files to a remote webdav server
type Webdav struct {
	client          *webdavClient.Webdav
	logger          *logrus.Logger
	UploadDirectory string
}

func init() {
	archivers.Register(New)
}

// New will return a new webdav uploader
func New(c interface{}, logger *logrus.Logger) archivers.Archiver {
	w := &Webdav{
		logger: logger,
	}
	config.ConfigFromStruct(c, &w)
	w.client = webdavClient.New(c, logger)
	return w
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(f file.File) (err error) {
	if err = w.Connect(); err != nil {
		return
	}

	uploadFilePath := path.Join(w.UploadDirectory, f.Directory)
	return w.client.Upload(f.Filename, uploadFilePath, f.Body)
}

func (w *Webdav) Connect() (err error) {
	return w.client.Connect()
}
