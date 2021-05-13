package file

import "io"

type File struct {
	Filename  string
	Directory string
	Body      io.Reader
}
