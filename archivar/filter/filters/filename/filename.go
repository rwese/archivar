package filename

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"

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
	for _, reject := range fc.Reject {
		f.reject = append(f.reject, regexp.MustCompile(reject))
	}

	for _, allow := range fc.Allow {
		f.allow = append(f.allow, regexp.MustCompile(allow))
	}

	return f
}

func (f *Filename) Filter(filename *string, filepath *string, data io.Reader) (filterResult.Results, error) {
	return f.runFilters(filename, filepath, data)
}

func (f *Filename) sanatizeName(filename *string) {
	cleanFilename := strings.TrimSpace(*filename)
	alphaRegex := regexp.MustCompile(`[^[:word:]-_. ]`)
	cleanFilename = alphaRegex.ReplaceAllString(cleanFilename, "")
	*filename = cleanFilename
}

func (f *Filename) runFilters(filename *string, filepath *string, data io.Reader) (filterResult.Results, error) {
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

	return filterResult.NoAction, nil
}
