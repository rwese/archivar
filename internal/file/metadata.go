package file

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"time"
)

var errNotFound = errors.New("key not set")
var errStoredValueNotParsable = errors.New("stored data is not usable as the requested type")
var errKeyIsReserved = errors.New("key is reserved")

type ChecksumFunc func(File) (checksum string)

const (
	MetaDataBody      = "Body"
	MetaDataFilename  = "Filename"
	MetaDataDirectory = "Directory"
	MetaDataChecksum  = "Checksum"
	MetaDataCreatedAt = "CreatedAt"
	MetaDataChangedAt = "ChangedAt"
)

var reservedKeys = map[string]bool{
	MetaDataBody:      true,
	MetaDataFilename:  true,
	MetaDataDirectory: true,
	MetaDataChecksum:  true,
	MetaDataCreatedAt: true,
	MetaDataChangedAt: true,
}

func (f *File) SetMetadataTime(key string, data time.Time) *File {
	f.Metadata[key] = data
	return f
}

func (f *File) getMetadataTime(key string) (stored time.Time, err error) {
	storedMetadata, ok := f.Metadata[key]
	if !ok {
		return time.Time{}, errNotFound
	}

	return storedMetadata.(time.Time), nil
}

func (f *File) setMetadataString(key string, data interface{}) {
	f.Metadata[key] = data
}

func (f *File) SetMetadataString(key string, data string) (err error) {
	if _, exists := reservedKeys[key]; exists {
		return errKeyIsReserved
	}

	f.setMetadata(key, data)
	return
}

func (f *File) setMetadata(key string, data interface{}) {
	f.Metadata[key] = data
}

func (f *File) GetMetadataString(key string) (data string, err error) {
	dataInterface, exists := f.Metadata[key]
	if !exists {
		return "", nil
	}

	data, ok := dataInterface.(string)
	if !ok {
		return "", errStoredValueNotParsable
	}

	return
}

func Checksum(f File) (checksum string) {
	return
}

func (f File) Checksum() (checksum string) {
	if f.ChecksumFunc == nil {
		return ""
	}

	return f.ChecksumFunc(f)
}

func ChecksumChanged(a, b File) bool {
	return a.Checksum() != b.Checksum()
}

type Metadata func(f *File)

func (f *File) SetMetadata(metadata ...Metadata) {
	for _, m := range metadata {
		m(f)
	}
}

func WithChecksumFunc(checksumFunc ChecksumFunc) Metadata {
	return func(f *File) {
		f.ChecksumFunc = checksumFunc
	}
}

func WithContent(content io.Reader) Metadata {
	return func(f *File) {
		f.Body = content
	}
}
func WithContentBase64Encoded(content io.Reader) Metadata {
	return func(f *File) {
		c, _ := io.ReadAll(content)
		decodedContent, _ := base64.StdEncoding.DecodeString(string(c))
		f.Body = io.MultiReader(bytes.NewReader(decodedContent))
	}
}

func WithContentPointer(content *io.Reader) Metadata {
	return func(f *File) {
		f.Body = *content
	}
}

func WithChangedAt(at time.Time) Metadata {
	return func(f *File) {
		f.SetMetadataTime(MetaDataChangedAt, at)
	}
}

func WithCreatedAt(createdAt time.Time) Metadata {
	return func(f *File) {
		f.SetMetadataTime(MetaDataCreatedAt, createdAt)
	}
}

func WithDirectory(dir string) Metadata {
	return func(f *File) {
		f.setMetadataString(MetaDataDirectory, dir)
	}
}

func WithFilename(filename string) Metadata {
	return func(f *File) {
		f.setMetadataString(MetaDataFilename, filename)
	}
}
