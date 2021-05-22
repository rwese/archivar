package webdav

import (
	"io"

	"github.com/rwese/archivar/internal/file"
)

// DownloadFile returns a io.Reader to retrieve the requested file
func (w *Webdav) DownloadFile(file string) (fileHandle io.Reader, err error) {

	_, err = w.Client.Stat(file)
	if err != nil {
		w.logger.Error("file does not exist")
	}

	return w.Client.ReadStream(file)
}

func (w *Webdav) DownloadFiles(directory string, fileCh chan file.File) (err error) {
	defer close(fileCh)

	files, err := w.Client.ReadDir(directory)
	if err != nil {
		w.logger.Error("file does not exist")
	}

	for _, f := range files {
		fh, err := w.DownloadFile(f.Name())
		if err != nil {
			w.logger.Warn(err)
		}

		fileCh <- file.File{
			Body:     fh,
			Filename: f.Name(),
		}
	}

	return
}

func (w *Webdav) DeleteFile(file file.File) (err error) {
	return w.Client.Remove(file.Path())
}

func (w *Webdav) DeleteFiles(files []file.File) (err error) {
	for _, file := range files {
		if err = w.DeleteFile(file); err != nil {
			return err
		}
	}

	return
}

// DownloadDirectory returns a io.Reader to retrieve the requested file
// func (w *Webdav) DownloadDirectory(directory string) (fileHandle io.Reader, err error) {
// 	_, err = w.Client.Stat(directory)
// 	if err != nil {
// 		w.logger.Error("directory does not exist")
// 	}

// 	return w.Client.ReadStream(file)
// }
