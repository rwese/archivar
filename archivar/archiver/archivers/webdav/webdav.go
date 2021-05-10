package webdav

import (
	"encoding/json"
	"io"
	"path"

	webdavClient "github.com/rwese/archivar/internal/webdav"

	"github.com/sirupsen/logrus"
)

// Webdav allows to upload files to a remote webdav server
type Webdav struct {
	client          *webdavClient.Webdav
	logger          *logrus.Logger
	UploadDirectory string
}

// New will return a new webdav uploader
func New(config interface{}, logger *logrus.Logger) (webdav *Webdav) {
	jsonM, _ := json.Marshal(config)
	json.Unmarshal(jsonM, &webdav)
	webdav.client = webdavClient.New(config, logger)
	return webdav
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	_, err = w.client.Connect()
	if err != nil {
		return
	}

	uploadFilePath := path.Join(w.UploadDirectory, fileDirectory)
	return w.client.Upload(fileName, uploadFilePath, fileHandle)
}
