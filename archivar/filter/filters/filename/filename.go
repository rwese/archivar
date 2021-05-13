package filename

import (
	"regexp"

	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/utils/config"

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

func New(c interface{}, logger *logrus.Logger) *Filename {
	var fc FilenameConfig
	config.ConfigFromStruct(c, &fc)
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

func (filename *Filename) Filter(f *file.File) (filterResult.Results, error) {
	for _, allow := range filename.allow {
		if allow.Match([]byte(f.Filename)) {
			filename.logger.Debugf("Filename: Allow %s", f.Filename)
			return filterResult.Allow, nil
		}
	}

	for _, reject := range filename.reject {
		if reject.Match([]byte(f.Filename)) {
			filename.logger.Debugf("Filename: Reject %s", f.Filename)
			return filterResult.Reject, nil
		}
	}

	filename.logger.Debugf("Filename: Miss %s", f.Filename)
	return filterResult.Miss, nil
}
