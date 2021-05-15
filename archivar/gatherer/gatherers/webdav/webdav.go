package webdav

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/gatherer/gatherers"
	"github.com/rwese/archivar/internal/file"
	webdavClient "github.com/rwese/archivar/internal/webdav"
	"github.com/rwese/archivar/utils/config"
	"github.com/sirupsen/logrus"
)

// Webdav allows to upload files to a remote webdav server
type Webdav struct {
	storage archivers.Archiver
	logger  *logrus.Logger
	client  *webdavClient.Webdav
}

func init() {
	gatherers.Register(New)
}

// New will return a new webdav downloader
func New(c interface{}, storage archivers.Archiver, logger *logrus.Logger) gatherers.Gatherer {
	webdav := &Webdav{
		storage: storage,
		logger:  logger,
		client:  webdavClient.New(c, logger),
	}
	config.ConfigFromStruct(c, &webdav)
	return webdav
}

func (w Webdav) Download() (err error) {
	if err = w.Connect(); err != nil {
		return
	}

	files := make(chan file.File)
	if err = w.client.DownloadFiles("", files); err != nil {
		return
	}

	for file := range files {
		w.storage.Upload(file)
	}

	return
}

func (w *Webdav) Connect() (err error) {
	if err = w.storage.Connect(); err != nil {
		return
	}

	return w.client.Connect()
}
