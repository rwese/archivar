package webdav

import (
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/internal/file"
	webdavClient "github.com/rwese/archivar/internal/webdav"
	"github.com/rwese/archivar/utils/config"
	"github.com/sirupsen/logrus"
)

// Webdav allows to upload files to a remote webdav server
type Webdav struct {
	storage archiver.Archiver
	logger  *logrus.Logger
	client  *webdavClient.Webdav
}

// New will return a new webdav downloader
func New(c interface{}, storage archiver.Archiver, logger *logrus.Logger) *Webdav {
	webdav := &Webdav{
		storage: storage,
		logger:  logger,
		client:  webdavClient.New(c, logger),
	}
	config.ConfigFromStruct(c, &webdav)
	return webdav
}

func (w *Webdav) Download() (err error) {
	files := make(chan file.File)
	if err = w.client.DownloadFiles("", files); err != nil {
		return
	}

	for file := range files {
		w.storage.Upload(file)
	}

	return
}
