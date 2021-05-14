package filename_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters/filename"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

func TestFilenameAccept(t *testing.T) {
	fileTests := map[string]struct {
		config  filename.FilenameConfig
		have    file.File
		want    file.File
		wantErr bool
		result  filterResult.Results
	}{
		"allow_first_do_not_modify_other_stuff": {
			config: filename.FilenameConfig{
				Allow: []string{
					"^allowme$",
				},
				Reject: []string{
					"^allowme$",
				},
			},
			have:   file.File{Filename: "allowme", Directory: "/somepath/", Body: bytes.NewReader([]byte(` Testing `))},
			want:   file.File{Filename: "allowme", Directory: "/somepath/", Body: bytes.NewReader([]byte(` Testing `))},
			result: filterResult.Allow,
		},
		"allow_reject": {
			config: filename.FilenameConfig{
				Allow: []string{
					"^allowme$",
				},
				Reject: []string{
					"^reject$",
				},
			},
			have:   file.File{Filename: "reject"},
			want:   file.File{Filename: "reject"},
			result: filterResult.Reject,
		},
		"reject": {
			config: filename.FilenameConfig{
				Reject: []string{
					"^reject$",
				},
			},
			have:   file.File{Filename: "reject"},
			want:   file.File{Filename: "reject"},
			result: filterResult.Reject,
		},
		"rejectPArtialRegex": {
			config: filename.FilenameConfig{
				Reject: []string{
					"reject",
				},
			},
			have:   file.File{Filename: "rejectThis"},
			want:   file.File{Filename: "rejectThis"},
			result: filterResult.Reject,
		},
	}

	for testName, fileTest := range fileTests {
		f := filename.New(fileTest.config, logrus.New())
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
