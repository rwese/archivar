package webdav

import (
	"fmt"
	"path"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/archiver/archivers/webdav/client"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
)

// NewArchiver will return a new webdav uploader
func NewArchiver(c interface{}, logger *logrus.Logger) archivers.Archiver {
	w := &WebdavArchiver{
		logger: logger,
	}
	config.ConfigFromStruct(c, &w)
	w.client = client.New(c, logger)
	return w
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *WebdavArchiver) Upload(f file.File) (err error) {
	if err = w.Connect(); err != nil {
		return
	}

	uploadFilePath := path.Join(w.UploadDirectory, f.Directory)
	uploadFile := path.Join(w.UploadDirectory, f.Directory, f.Filename)

	if w.compareChecksum(uploadFile, f.Checksum) {
		return nil
	}

	return w.client.Upload(f.Filename, uploadFilePath, f.Body)
}

func (w *WebdavArchiver) Connect() (err error) {
	return w.client.Connect()
}

func (w *WebdavArchiver) compareChecksum(file, checksum string) bool {
	if checksum == "" {
		return false
	}

	fs, err := w.client.Client.Stat(file)
	if err != nil {
		return false
	}

	currentChecksum := fmt.Sprintf("%d", fs.Size())
	return checksum == currentChecksum
}
