package filesystem

import (
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/rwese/archivar/internal/file"
)

type FileSystem struct{}

func New(directory string) *FileSystem {
	return &FileSystem{}
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

func (f *FileSystem) DownloadFile(file file.File) ([]byte, error) {
	return os.ReadFile(file.Path())
}

func (f *FileSystem) ListFiles(directory string) ([]file.File, error) {
	dir, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var files []file.File
	for _, fentry := range dir {
		if fentry.IsDir() {
			continue
		}

		files = append(files, file.File{
			Filename:  fentry.Name(),
			Directory: directory,
		})
	}

	return files, nil
}
