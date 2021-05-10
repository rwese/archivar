package webdav

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

// Webdav allows to upload files to a remote webdav server
type Webdav struct {
	Server                 string
	UserName               string
	Password               string
	UploadDirectory        string
	isRetry                bool
	knownUploadDirectories map[string]bool
	logger                 *logrus.Logger
	client                 *gowebdav.Client
}

// New will return a new webdav uploader
func New(config interface{}, logger *logrus.Logger) *Webdav {
	webdav := &Webdav{
		logger: logger,
	}
	webdav.knownUploadDirectories = make(map[string]bool)
	jsonM, _ := json.Marshal(config)
	json.Unmarshal(jsonM, &webdav)
	return webdav
}

func (w *Webdav) Connect() (err error) {

	return
}

func (w *Webdav) Download() (err error) {

	return
}
