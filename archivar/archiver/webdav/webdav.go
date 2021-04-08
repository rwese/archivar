package webdav

import (
	"io"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

// Webdav allows to upload files to a remote webdav server
type Webdav struct {
	webRoot         string
	userName        string
	userPassword    string
	keepUploaded    bool
	uploadDirectory string
	logger          *logrus.Logger
}

// New will return a new webdav uploader
func New(webRoot string, userName string, userPassword string, uploadDirectory string, logger *logrus.Logger) Webdav {
	return Webdav{
		webRoot:         webRoot,
		userName:        userName,
		userPassword:    userPassword,
		keepUploaded:    false,
		uploadDirectory: uploadDirectory,
		logger:          logger,
	}
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w Webdav) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	c := gowebdav.NewClient(w.webRoot, w.userName, w.userPassword)

	w.logger.Debugf("connecting to %s as %s\n", w.webRoot, w.userName)
	err = c.Connect()
	if err != nil {
		w.logger.Fatalf("failed to connect: %s\n", err.Error())
	}

	_, err = c.Stat(w.uploadDirectory)
	if err != nil {
		w.logger.Fatalf("failed to access upload directory: %s", err.Error())
	}

	uploadDirectory := path.Join(w.uploadDirectory, fileDirectory)
	w.logger.Debugf("uploadDirectory will be: %s", uploadDirectory)
	_, err = c.Stat(uploadDirectory)
	if err != nil {
		err = c.MkdirAll(uploadDirectory, 0644)
		if err != nil {
			w.logger.Fatalf("failed to create uploadTargetFolder: %s", err.Error())
			return
		}
	}

	uploadFileName := path.Join(uploadDirectory, fileName)
	return c.WriteStream(uploadFileName, fileHandle, 0644)
}
