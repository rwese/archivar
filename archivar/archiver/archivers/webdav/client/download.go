package client

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/rwese/archivar/archivar/archiver/archivers"
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

func (w *Webdav) DownloadFiles(directory string, upload archivers.UploadFunc) ([]string, error) {
	return w.downloadFilesRecursive(directory, directory, upload)
}

func (w *Webdav) downloadFilesRecursive(rootDirectory, directory string, upload archivers.UploadFunc) (downloadedFiles []string, err error) {
	files, err := w.Client.ReadDir(directory)
	if err != nil {
		w.logger.Error("file does not exist")
	}

	for _, f := range files {
		fullPath := path.Join(directory, f.Name())
		if f.IsDir() {
			rdownloadedFiles, err := w.downloadFilesRecursive(rootDirectory, fullPath, upload)
			if err != nil {
				return nil, err
			}

			downloadedFiles = append(downloadedFiles, rdownloadedFiles...)
			continue
		}

		if w.SeenFiles[fullPath] != 0 && w.SeenFiles[fullPath] == f.ModTime().Unix() {
			continue
		}

		fh, err := w.DownloadFile(fullPath)
		if err != nil {
			w.logger.Warn(err)
		}

		relativeDirectory := strings.TrimPrefix(directory, rootDirectory)
		newFile := file.New(
			file.WithContent(fh),
			file.WithFilename(f.Name()),
			file.WithDirectory(relativeDirectory),
			nil,
		)

		newFile.SetMetadataString("Filesize", fmt.Sprintf("%d", f.Size()))
		newFile.ChecksumFunc = func(f file.File) string {
			data, err := f.GetMetadataString("Filesize")
			if err != nil {
				return ""
			}

			return data
		}

		if err = upload(newFile); err != nil {
			return nil, err
		}

		w.SeenFiles[fullPath] = f.ModTime().Unix()
		downloadedFiles = append(downloadedFiles, fullPath)
	}

	return
}

func (w *Webdav) DeleteFile(file string) (err error) {
	return w.Client.Remove(file)
}

func (w *Webdav) DeleteFiles(files []string) (err error) {
	for _, file := range files {
		if err = w.DeleteFile(file); err != nil {
			return err
		}
	}

	return
}
