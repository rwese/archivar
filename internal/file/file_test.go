package file_test

import (
	"bytes"
	"testing"

	"github.com/rwese/archivar/internal/file"
)

func TestFile(t *testing.T) {
	testFilename := "Filename"
	testDirectory := "Directory"
	file := file.New(testFilename, testDirectory, bytes.NewReader([]byte("body")), nil)
	if file.Filename() != testFilename {
		t.Fatal("Filename missmatch")
	}

	if file.Directory() != testDirectory {
		t.Fatal("Directory missmatch")
	}
}
