package webdav

import (
	"os"

	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

type Webdav struct {
	Server           string
	UserName         string
	Password         string
	KnownDirectories map[string]bool  // KnownDirectories improve speed by reducing repeated directory lookups
	SeenFiles        map[string]int64 // Already seen files won't re-reuploaded
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
	webdav.SeenFiles = make(map[string]int64)
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

func (w *Webdav) DirExists(d string) bool {
	if w.KnownDirectories[d] {
		return true
	}

	ds, err := w.Client.Stat(d)

	if _, ok := err.(*os.PathError); ok {
		return false
	}

	w.KnownDirectories[d] = true

	if !ds.IsDir() {
		w.logger.Warnf("Directory '%s' is no directory.", d)
		return true
	}

	return true
}
