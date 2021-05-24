package file

import (
	"errors"
)

var errStoredValueNotParsable = errors.New("stored data is not usable as the requested type")
var errKeyIsReserved = errors.New("key is reserved")

type ChecksumFunc func(File) (checksum string)

const MetaDataFilename = "Filename"
const MetaDataDirectory = "Directory"
const MetaDataChecksum = "Checksum"

var reservedKeys = map[string]bool{
	MetaDataFilename:  true,
	MetaDataDirectory: true,
	MetaDataChecksum:  true,
}

func (f *File) SetMetadata(key string, data interface{}) (err error) {
	if _, exists := reservedKeys[key]; exists {
		return errKeyIsReserved
	}

	f.setMetadataString(key, data)
	return
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

func (f *File) SetFilename(Filename string) {
	f.setMetadataString(MetaDataFilename, Filename)
}

func (f *File) SetDirectory(Directory string) {
	f.setMetadataString(MetaDataDirectory, Directory)
}

func (f *File) Filename() string {
	data, err := f.GetMetadataString(MetaDataFilename)
	if err != nil {
		return ""
	}

	return data
}

func (f *File) Directory() string {
	data, err := f.GetMetadataString(MetaDataDirectory)
	if err != nil {
		return ""
	}

	return data
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

func FileChanged(a, b File) bool {
	return a.Checksum() != b.Checksum()
}
