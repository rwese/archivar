package webdav

import (
	"io"
	"path"

	"github.com/studio-b12/gowebdav"
)

var knownUploadDirectories map[string]bool

func (w *Webdav) connect() (newSession bool, err error) {
	if newSession = w.client == nil; !newSession {
		return
	}

	w.client = gowebdav.NewClient(w.server, w.userName, w.userPassword)
	w.logger.Debugf("connecting to %s as %s\n", w.server, w.userName)
	err = w.client.Connect()
	if err != nil {
		w.logger.Fatalf("failed to connect: %s\n", err.Error())
	}

	return
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {

	newSession, err := w.connect()
	if err != nil {
		return
	}

	if newSession {
		knownUploadDirectories = make(map[string]bool)
		_, err = w.client.Stat(w.uploadDirectory)
		if err != nil {
			w.logger.Fatalf("failed to access upload directory: %s", err.Error())
		}
	}

	uploadDirectory := path.Join(w.uploadDirectory, fileDirectory)

	if !knownUploadDirectories[uploadDirectory] {
		w.logger.Debugf("uploadDirectory will be: %s", uploadDirectory)
		_, err = w.client.Stat(uploadDirectory)
		if err != nil {
			err = w.client.MkdirAll(uploadDirectory, 0644)
			if err != nil {
				w.logger.Fatalf("failed to create uploadTargetFolder: %s", err.Error())
				return
			}
		}
		knownUploadDirectories[uploadDirectory] = true
	}

	uploadFileName := path.Join(uploadDirectory, fileName)
	w.logger.Debugf("uploading: %s", uploadFileName)
	return w.client.WriteStream(uploadFileName, fileHandle, 0644)
}
