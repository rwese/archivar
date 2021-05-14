package sanatizer_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/rwese/archivar/archivar/processor/processors/sanatizer"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

func TestSanatizeTrim(t *testing.T) {
	fileTests := map[string]struct {
		config sanatizer.SanatizeConfig
		have   file.File
		want   file.File
	}{
		"only trim filename": {
			config: sanatizer.SanatizeConfig{TrimWhitespaces: true},
			have:   file.File{Filename: " whitespace_before", Directory: "/somepath/ ", Body: bytes.NewReader([]byte(` Testing `))},
			want:   file.File{Filename: "whitespace_before", Directory: "/somepath/ ", Body: bytes.NewReader([]byte(` Testing `))},
		},
	}

	for testName, fileTest := range fileTests {
		f := sanatizer.New(fileTest.config, logrus.New())

		file := fileTest.have
		f.Process(&file)
		if !reflect.DeepEqual(file, fileTest.want) {
			t.Fatalf("Failed test '%s'", testName)
		}
	}
}
func TestSanatizeCharacterBlacklistRegexs(t *testing.T) {
	fileTests := map[string]struct {
		config sanatizer.SanatizeConfig
		have   file.File
		want   file.File
	}{
		"simple replace": {
			config: sanatizer.SanatizeConfig{CharacterBlacklistRegexs: []string{"[0-9]"}},
			have:   file.File{Filename: "a1b2", Directory: "/somepath/ ", Body: bytes.NewReader([]byte(` Testing `))},
			want:   file.File{Filename: "ab", Directory: "/somepath/ ", Body: bytes.NewReader([]byte(` Testing `))},
		},
	}

	for testName, fileTest := range fileTests {
		f := sanatizer.New(fileTest.config, logrus.New())

		file := fileTest.have
		f.Process(&file)
		if !reflect.DeepEqual(file, fileTest.want) {
			t.Fatalf("Failed test '%s'", testName)
		}
	}
}
