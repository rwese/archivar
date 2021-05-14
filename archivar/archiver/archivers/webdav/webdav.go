package webdav

import (
	"path"

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

// New will return a new webdav uploader
func New(c interface{}, logger *logrus.Logger) (webdav *Webdav) {
	config.ConfigFromStruct(c, &webdav)
	webdav.client = webdavClient.New(c, logger)
	return webdav
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(f file.File) (err error) {
	_, err = w.client.Connect()
	if err != nil {
		return
	}

	uploadFilePath := path.Join(w.UploadDirectory, f.Directory)
	return w.client.Upload(f.Filename, uploadFilePath, f.Body)
}
