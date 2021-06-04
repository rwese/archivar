package file_test

import (
	"bytes"
	"testing"

	"github.com/rwese/archivar/internal/file"
)

func TestFile(t *testing.T) {
	testFilename := "Filename"
	testDirectory := "Directory"
	file := file.New(
		file.WithContent(bytes.NewReader([]byte("body"))),
		file.WithDirectory(testDirectory),
		file.WithFilename(testFilename),
	)
	if file.Filename() != testFilename {
		t.Fatal("Filename missmatch")
	}

	if file.Directory() != testDirectory {
		t.Fatal("Directory missmatch")
	}
}
