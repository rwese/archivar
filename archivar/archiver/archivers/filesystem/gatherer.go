package filesystem

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/archiver/archivers/filesystem/client"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
)

//
type FileSystemGatherer struct {
	storage          archivers.Archiver
	logger           *logrus.Logger
	client           *client.FileSystem
	directory        string
	deleteDownloaded bool
}

type FileSystemGathererConfig struct {
	Directory        string
	DeleteDownloaded bool
}

func NewGatherer(c interface{}, storage archivers.Archiver, logger *logrus.Logger) archivers.Gatherer {
	wc := &FileSystemGathererConfig{}
	config.ConfigFromStruct(c, &wc)

	filesystem := &FileSystemGatherer{
		storage:          storage,
		logger:           logger,
		client:           client.New(logger),
		directory:        wc.Directory,
		deleteDownloaded: wc.DeleteDownloaded,
	}
	return filesystem
}

func (w FileSystemGatherer) Download() (err error) {
	if err = w.Connect(); err != nil {
		return
	}

	var downloadedFiles []string
	if err = w.client.DownloadFiles(w.directory, w.storage.Upload); err != nil {
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

func (w *FileSystemGatherer) Connect() (err error) {
	if !w.client.DirExists(w.directory) {
		w.logger.Fatalf("failed to access upload directory, which will not be automatically created: %s", err.Error())
	}
	return
}
