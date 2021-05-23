package file

import (
	"io"
	"os"
	"path"
)

type File struct {
	Filename  string
	Directory string
	Checksum  string // Checksum is the filesize in bytes, not the best, but works for now
	Body      io.Reader
}

func (f File) FullFilePath() string {
	return path.Join(f.Directory, f.Filename)
}

func (f File) Exists() bool {
	if _, err := os.Stat(f.FullFilePath()); os.IsNotExist(err) {
		return false
	}

	return true
}
