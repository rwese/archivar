package filename

import (
	"encoding/json"
	"io"
	"regexp"

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
	f := &Filename{
		logger: logger,
	}

	jsonM, _ := json.Marshal(config)
	var fc FilenameConfig
	json.Unmarshal(jsonM, &fc)
	for _, reject := range fc.Reject {
		f.reject = append(f.reject, regexp.MustCompile(reject))
	}

	for _, allow := range fc.Allow {
		f.allow = append(f.allow, regexp.MustCompile(allow))
	}

	return f
}

func (f *Filename) Filter(filename, filepath string, data io.Reader) (bool, error) {
	for _, allow := range f.allow {
		if allow.Match([]byte(filename)) {
			f.logger.Debugf("Allow %s", filename)
			return true, nil
		}
	}

	for _, reject := range f.reject {
		if reject.Match([]byte(filename)) {
			f.logger.Debugf("Reject %s", filename)
			return false, nil
		}
	}

	return true, nil
}
