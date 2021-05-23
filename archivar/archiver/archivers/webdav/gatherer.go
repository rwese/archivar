package webdav

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/archiver/archivers/webdav/client"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
)

// WebdavGatherer allows to upload files to a remote webdav server
type WebdavGatherer struct {
	storage          archivers.Archiver
	logger           *logrus.Logger
	client           *client.Webdav
	directory        string
	deleteDownloaded bool
}

type WebdavGathererConfig struct {
	Directory        string
	DeleteDownloaded bool
}

// New will return a new webdav downloader
func NewGatherer(c interface{}, storage archivers.Archiver, logger *logrus.Logger) archivers.Gatherer {
	wc := &WebdavGathererConfig{}
	config.ConfigFromStruct(c, &wc)

	webdav := &WebdavGatherer{
		storage:          storage,
		logger:           logger,
		client:           client.New(c, logger),
		directory:        wc.Directory,
		deleteDownloaded: wc.DeleteDownloaded,
	}
	return webdav
}

func (w WebdavGatherer) Download() (err error) {
	if err = w.Connect(); err != nil {
		return
	}

	var downloadedFiles []string
	if downloadedFiles, err = w.client.DownloadFiles(w.directory, w.storage.Upload); err != nil {
		return
	}

	if w.deleteDownloaded {
		err = w.client.DeleteFiles(downloadedFiles)
		if err != nil {
			return err
		}
	}

	return
}

func (w *WebdavGatherer) Connect() (err error) {
	if err = w.client.Connect(); err != nil {
		return
	}

	if !w.client.DirExists(w.directory) {
		w.logger.Fatalf("failed to access upload directory, which will not be automatically created: %s", err.Error())
	}
	return
}
