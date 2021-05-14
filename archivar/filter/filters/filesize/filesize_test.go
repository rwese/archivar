package filesize_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters/filesize"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

func TestFilesizeMin(t *testing.T) {
	fileTests := map[string]struct {
		config  filesize.FilesizeConfig
		have    file.File
		want    file.File
		wantErr bool
		result  filterResult.Results
	}{
		"minSize_ok": {
			config: filesize.FilesizeConfig{MinSizeBytes: 10},
			have:   file.File{Body: bytes.NewReader([]byte(`exactly10.`))},
			want:   file.File{Body: bytes.NewReader([]byte(`exactly10.`))},
			result: filterResult.Allow,
		},
		"minSize_nok": {
			config: filesize.FilesizeConfig{MinSizeBytes: 10},
			have:   file.File{Body: bytes.NewReader([]byte(`oneshort!`))},
			want:   file.File{Body: bytes.NewReader([]byte(`oneshort!`))},
			result: filterResult.Reject,
		},
		"maxSize_ok": {
			config: filesize.FilesizeConfig{MaxSizeBytes: 10},
			have:   file.File{Body: bytes.NewReader([]byte(`exactly10.`))},
			want:   file.File{Body: bytes.NewReader([]byte(`exactly10.`))},
			result: filterResult.Allow,
		},
		"maxSize_nok": {
			config: filesize.FilesizeConfig{MaxSizeBytes: 10},
			have:   file.File{Body: bytes.NewReader([]byte(`justoneover`))},
			want:   file.File{Body: bytes.NewReader([]byte(`justoneover`))},
			result: filterResult.Reject,
		},
	}

	for testName, fileTest := range fileTests {
		f := filesize.New(fileTest.config, logrus.New())
		file := fileTest.have
		result, err := f.Filter(&file)

		if fileTest.result != result {
			t.Fatalf("'%s' wrong result", testName)
		}

		if fileTest.wantErr && err == nil {
			t.Fatalf("'%s' wantErr", testName)
		}

		if !reflect.DeepEqual(file, fileTest.want) && !fileTest.wantErr && result != filterResult.Reject {
			t.Fatalf("'%s' Failed test missmatch", testName)
		}
	}
}
