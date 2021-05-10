package filename

import (
	"encoding/json"
	"io"
	"regexp"

	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/sirupsen/logrus"
)

type FilenameConfig struct {
	Reject []string
	Allow  []string
}

type Filename struct {
	reject []*regexp.Regexp
	allow  []*regexp.Regexp
	logger *logrus.Logger
}

func New(config interface{}, logger *logrus.Logger) *Filename {
	jsonM, _ := json.Marshal(config)
	var fc FilenameConfig
	json.Unmarshal(jsonM, &fc)

	f := &Filename{
		logger: logger,
	}

	if len(fc.Reject) == 0 && len(fc.Allow) == 0 {
		logger.Fatalln("Filename filter requires at least one reject/allow rule")
	}

	for _, reject := range fc.Reject {
		f.reject = append(f.reject, regexp.MustCompile(reject))
	}

	for _, allow := range fc.Allow {
		f.allow = append(f.allow, regexp.MustCompile(allow))
	}

	return f
}

func (f *Filename) Filter(filename *string, filepath *string, data *io.Reader) (filterResult.Results, error) {
	for _, allow := range f.allow {
		if allow.Match([]byte(*filename)) {
			f.logger.Debugf("Allow %s", *filename)
			return filterResult.Allow, nil
		}
	}

	for _, reject := range f.reject {
		if reject.Match([]byte(*filename)) {
			f.logger.Debugf("Reject %s", *filename)
			return filterResult.Reject, nil
		}
	}

	return filterResult.Miss, nil
}
