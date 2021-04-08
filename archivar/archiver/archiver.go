package archiver

import "io"

type Archiver interface {
	Upload(fileName string, directory string, fileHandle io.Reader) (err error)
}
