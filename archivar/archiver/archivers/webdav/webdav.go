package webdav

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/archiver/archivers/webdav/client"

	"github.com/sirupsen/logrus"
)

// WebdavArchiver allows to upload files to a remote webdav server
type WebdavArchiver struct {
	client *client.Webdav

	logger          *logrus.Logger
	UploadDirectory string
}

func init() {
	archivers.RegisterArchiver(NewArchiver)
	archivers.RegisterGatherer(NewGatherer)
}
