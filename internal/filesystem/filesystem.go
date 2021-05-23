package filesystem

import (
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

type FileSystem struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *FileSystem {
	return &FileSystem{
		logger: logger,
	}
}

func (f *FileSystem) Upload(
	fileName string,
	fileDirectory string,
	fileHandle io.Reader,
) (err error) {
	if err = f.MkdirAll(fileDirectory); err != nil {
		return
	}

	filePath := path.Join(fileDirectory, fileName)
	fileData, err := io.ReadAll(fileHandle)
	if err != nil {
		return
	}
	return os.WriteFile(filePath, fileData, fs.FileMode(0660))
}

func (f *FileSystem) MkdirAll(d string) error {
	return os.MkdirAll(d, 0760)
}

func (f *FileSystem) DownloadFiles(d string, upload archivers.UploadFunc) (err error) {
	files, err := f.ListFiles(d)
	if err != nil {
		return
	}

	for _, file := range files {
		if err = f.DownloadFile(file, upload); err != nil {
			return
		}
	}

	return
}

func (f *FileSystem) DownloadFile(file file.File, upload archivers.UploadFunc) error {
	return upload(file)
}

func (f *FileSystem) DirExists(d string) bool {
	_, err := os.Stat(d)
	return !os.IsNotExist(err)
}

func (f *FileSystem) DeleteFiles(files []string) (err error) {
	for _, file := range files {
		if err = os.Remove(file); err != nil {
			return
		}
	}
	return
}

func (f *FileSystem) ListFiles(directory string) ([]file.File, error) {
	return f.listFilesRecursive(directory, directory)
}

func (f *FileSystem) listFilesRecursive(rootdirectory, directory string) ([]file.File, error) {
	dir, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var files []file.File
	for _, fentry := range dir {
		fullPath := path.Join(directory, fentry.Name())

		if fentry.IsDir() {
			rfiles, err := f.listFilesRecursive(rootdirectory, fullPath)
			if err != nil {
				return nil, err
			}

			files = append(files, rfiles...)
			continue
		}

		fh, err := os.Open(fullPath)
		if err != nil {
			return nil, err
		}

		cleanDirectory := strings.TrimSuffix(directory, rootdirectory)
		files = append(files, file.File{
			Filename:  fentry.Name(),
			Directory: cleanDirectory,
			Body:      fh,
		})
	}

	return files, nil
}
