package webdav

import (
	"errors"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"
)

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	newSession, err := w.Connect()
	if err != nil {
		return
	}

	if newSession {
		_, err = w.client.Stat(w.UploadDirectory)
		if err != nil {
			w.logger.Fatalf("failed to access upload directory: %s", err.Error())
		}
	}

	uploadDirectory := path.Join(w.UploadDirectory, fileDirectory)

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

	err = w.client.WriteStream(uploadFileName, fileHandle, 0644)
	if !w.isRetry && isConflictError(err) {
		w.logger.Warnf("collision writing to: %s, retrying", uploadFileName)
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
		w.isRetry = true
		return w.Upload(fileName, fileDirectory, fileHandle)
	}
	w.isRetry = false
	return
}

func isConflictError(err error) bool {
	var perr *fs.PathError
	return (err != nil && errors.As(err, &perr) && perr.Unwrap().Error() == strconv.Itoa(http.StatusLocked))
}
