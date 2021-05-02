package webdav

import (
	"io"
	"path"
)

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	newSession, err := w.Connect()
	if err != nil {
		return
	}

	if newSession {
		_, err = w.client.Stat(w.uploadDirectory)
		if err != nil {
			w.logger.Fatalf("failed to access upload directory: %s", err.Error())
		}
	}

	uploadDirectory := path.Join(w.uploadDirectory, fileDirectory)

	if !w.knownUploadDirectories[uploadDirectory] {
		w.logger.Debugf("uploadDirectory will be: %s", uploadDirectory)
		_, err = w.client.Stat(uploadDirectory)
		if err != nil {
			err = w.client.MkdirAll(uploadDirectory, 0644)
			if err != nil {
				w.logger.Fatalf("failed to create uploadTargetFolder: %s", err.Error())
				return
			}
		}
		w.knownUploadDirectories[uploadDirectory] = true
	}

	uploadFileName := path.Join(uploadDirectory, fileName)
	w.logger.Debugf("uploading: %s", uploadFileName)
	return w.client.WriteStream(uploadFileName, fileHandle, 0644)
}
