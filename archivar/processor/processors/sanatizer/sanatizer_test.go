package sanatizer_test

import (
	"bytes"
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
			have: *file.New(
				file.WithContent(bytes.NewReader([]byte(` Testing `))),
				file.WithFilename("a1b2 "),
				file.WithDirectory("/somepath/ "),
			),
			want: *file.New(
				file.WithContent(bytes.NewReader([]byte(` Testing `))),
				file.WithFilename("a1b2"),
				file.WithDirectory("/somepath/ "),
			),
		},
	}

	for testName, fileTest := range fileTests {
		f := sanatizer.New(fileTest.config, logrus.New())

		file := fileTest.have
		f.Process(&file)
		if file.Filename() != fileTest.want.Filename() {
			t.Fatalf("Failed test '%s'", testName)
		}
		if file.Directory() != fileTest.want.Directory() {
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
			have: *file.New(
				file.WithContent(bytes.NewReader([]byte(` Testing `))),
				file.WithFilename("a1b2"),
				file.WithDirectory("/somepath/ "),
			),
			want: *file.New(
				file.WithContent(bytes.NewReader([]byte(` Testing `))),
				file.WithFilename("ab"),
				file.WithDirectory("/somepath/ "),
			),
		},
	}

	for testName, fileTest := range fileTests {
		f := sanatizer.New(fileTest.config, logrus.New())

		file := fileTest.have
		f.Process(&file)
		if file.Filename() != fileTest.want.Filename() {
			t.Fatalf("Failed test '%s'", testName)
		}
	}
}
