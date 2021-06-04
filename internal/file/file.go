package file

import (
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type File struct {
	Body         io.Reader
	Metadata     map[string]interface{}
	ChecksumFunc ChecksumFunc
}

func New(metadata ...Metadata) *File {
	f := &File{
		Metadata:     make(map[string]interface{}),
		ChecksumFunc: Checksum,
	}

	for _, m := range metadata {
		m(f)
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

func (f *File) Filename() string {
	data, err := f.GetMetadataString(MetaDataFilename)
	if err != nil {
		return ""
	}

	return data
}

func (f *File) SetFilename(Filename string) {
	Filename = strings.ReplaceAll(Filename, "/", "_")
	f.setMetadataString(MetaDataFilename, Filename)
}

func (f *File) SetDirectory(Directory string) {
	f.setMetadataString(MetaDataDirectory, Directory)
}

func (f *File) Directory() string {
	data, err := f.GetMetadataString(MetaDataDirectory)
	if err != nil {
		return ""
	}

	return data
}

func (f *File) SetChangedAt(t time.Time) {
	f.SetMetadataTime(MetaDataCreatedAt, t)
}

func (f *File) ChangedAt() (time.Time, error) {
	return f.getMetadataTime(MetaDataCreatedAt)
}
