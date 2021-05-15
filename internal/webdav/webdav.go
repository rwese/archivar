package webdav

import (
	"github.com/rwese/archivar/utils/config"
	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

type Webdav struct {
	Server           string
	UserName         string
	Password         string
	KnownDirectories map[string]bool // KnownDirectories improve speed by reducing repeated directory lookups
	isRetry          bool
	newSession       bool
	logger           *logrus.Logger
	Client           *gowebdav.Client
}

// New will return a new webdav uploader
func New(c interface{}, logger *logrus.Logger) *Webdav {
	webdav := &Webdav{logger: logger, newSession: true}
	config.ConfigFromStruct(c, &webdav)
	webdav.KnownDirectories = make(map[string]bool)
	return webdav
}

func (w *Webdav) Connect() (err error) {
	if w.Client != nil {
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
