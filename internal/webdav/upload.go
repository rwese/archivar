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

func (w *Webdav) createDirectoryIfNotExists(fileDirectory string) (err error) {
	if !w.KnownDirectories[fileDirectory] {
		w.logger.Debugf("fileDirectory will be: %s", fileDirectory)
		_, err = w.Client.Stat(fileDirectory)
		if err != nil {
			err = w.Client.MkdirAll(fileDirectory, 0644)
			if err != nil {
				w.logger.Fatalf("failed to create uploadTargetFolder: %s", err.Error())
				return
			}
		}
		w.KnownDirectories[fileDirectory] = true
	}

	return
}

// Upload takes filename, fileDirectory and fileHandle to push the data directly to the webdav
func (w *Webdav) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	if err = w.Connect(); err != nil {
		return
	}

	if w.newSession {
		w.newSession = false
		_, err = w.Client.Stat(fileDirectory)
		if err != nil {
			w.logger.Fatalf("failed to access upload directory, which will not automatically created: %s", err.Error())
		}

	}

	w.createDirectoryIfNotExists(fileDirectory)

	uploadFileName := path.Join(fileDirectory, fileName)
	w.logger.Debugf("uploading to: %s", uploadFileName)

	err = w.Client.WriteStream(uploadFileName, fileHandle, 0644)
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
