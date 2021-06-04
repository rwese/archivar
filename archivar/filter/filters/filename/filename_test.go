package filename_test

import (
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
			have: *file.New(
				file.WithFilename("allowme"),
				file.WithDirectory("/somepath/"),
			),
			want: *file.New(
				file.WithFilename("allowme"),
				file.WithDirectory("/somepath/"),
			),
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
			have: *file.New(
				file.WithFilename("reject"),
				file.WithDirectory(""),
			),
			want: *file.New(
				file.WithFilename("reject"),
				file.WithDirectory(""),
			),
			result: filterResult.Reject,
		},
		"reject": {
			config: filename.FilenameConfig{
				Reject: []string{
					"^reject$",
				},
			},
			have: *file.New(
				file.WithFilename("reject"),
				file.WithDirectory(""),
			),
			want: *file.New(
				file.WithFilename("reject"),
				file.WithDirectory(""),
			),
			result: filterResult.Reject,
		},
		"rejectPartialRegex": {
			config: filename.FilenameConfig{
				Reject: []string{
					"reject",
				},
			},
			have: *file.New(
				file.WithFilename("rejectThis"),
				file.WithDirectory(""),
			),
			want: *file.New(
				file.WithFilename("rejectThis"),
				file.WithDirectory(""),
			),
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
	}
}
