package file

import (
	"io"
	"os"
	"path"
)

type File struct {
	Body         io.Reader
	Metadata     map[string]interface{}
	ChecksumFunc ChecksumFunc
}

func New(filename, directory string, body io.Reader, cf ChecksumFunc) File {
	f := File{
		Body:         body,
		Metadata:     make(map[string]interface{}),
		ChecksumFunc: Checksum,
	}

	f.SetFilename(filename)
	f.SetDirectory(directory)

	if cf != nil {
		f.ChecksumFunc = cf
	}

	return f
}

func (f File) FullFilePath() string {
	return path.Join(f.Directory(), f.Filename())
}

func (f File) Exists() bool {
	if _, err := os.Stat(f.FullFilePath()); os.IsNotExist(err) {
		return false
	}

	return true
}
