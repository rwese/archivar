package filesize_test

import (
	"bytes"
	"io"
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
			have:   *file.New(file.WithContent(bytes.NewReader([]byte(`exactly10.`)))),
			want:   *file.New(file.WithContent(bytes.NewReader([]byte(`exactly10.`)))),
			result: filterResult.Allow,
		},
		"minSize_nok": {
			config: filesize.FilesizeConfig{MinSizeBytes: 10},
			have:   *file.New(file.WithContent(bytes.NewReader([]byte(`oneshort!`)))),
			want:   *file.New(file.WithContent(bytes.NewReader([]byte(`oneshort!`)))),
			result: filterResult.Reject,
		},
		"maxSize_nok": {
			config: filesize.FilesizeConfig{MaxSizeBytes: 10},
			have:   *file.New(file.WithContent(bytes.NewReader([]byte(`justoneover`)))),
			want:   *file.New(file.WithContent(bytes.NewReader([]byte(`justoneover`)))),
			result: filterResult.Reject,
		},
		"maxSize_ok": {
			config: filesize.FilesizeConfig{MaxSizeBytes: 10},
			have:   *file.New(file.WithContent(bytes.NewReader([]byte(`exactly10.`)))),
			want:   *file.New(file.WithContent(bytes.NewReader([]byte(`exactly10.`)))),
			result: filterResult.Allow,
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

		haveBuffer, _ := io.ReadAll(file.Body)
		wantBuffer, _ := io.ReadAll(fileTest.want.Body)
		if !bytes.Equal(haveBuffer, wantBuffer) && !fileTest.wantErr && result != filterResult.Reject {
			t.Fatalf("'%s' Failed test missmatch", testName)
		}
	}
}
