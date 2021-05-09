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
	logger                 *logrus.Logger
	knownUploadDirectories map[string]bool
	client                 *gowebdav.Client
	isRetry                bool
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
