package file

import (
	"io"
	"os"
	"path"
)

type File struct {
	Filename  string
	Directory string
	Body      io.Reader
}

func (f File) Path() string {
	return path.Join(f.Directory, f.Filename)
}

func (f File) Exists() bool {
	if _, err := os.Stat(f.Path()); os.IsNotExist(err) {
		return false
	}

	return true
}
