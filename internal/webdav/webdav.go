package webdav

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

type Webdav struct {
	Server                 string
	UserName               string
	Password               string
	KnownUploadDirectories map[string]bool
	isRetry                bool
	logger                 *logrus.Logger
	Client                 *gowebdav.Client
}

// New will return a new webdav uploader
func New(config interface{}, logger *logrus.Logger) *Webdav {
	webdav := &Webdav{logger: logger}
	jsonM, _ := json.Marshal(config)
	webdav.KnownUploadDirectories = make(map[string]bool)
	json.Unmarshal(jsonM, &webdav)
	return webdav
}

func (w *Webdav) Connect() (newSession bool, err error) {
	if newSession = w.Client == nil; !newSession {
		return
	}

	w.Client = gowebdav.NewClient(w.Server, w.UserName, w.Password)
	w.logger.Debugf("connecting to %s as %s\n", w.Server, w.UserName)
	err = w.Client.Connect()
	if err != nil {
		w.logger.Fatalf("failed to connect: %s\n", err.Error())
	}

	return
}
