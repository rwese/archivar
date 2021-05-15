package sanatizer

import (
	"regexp"
	"strings"

	"github.com/rwese/archivar/archivar/processor/processors"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/utils/config"
	"github.com/sirupsen/logrus"
)

type SanatizeConfig struct {
	TrimWhitespaces          bool
	CharacterBlacklistRegexs []string
}

type Sanatize struct {
	trimWhitespaces        bool
	logger                 *logrus.Logger
	characterReplaceRegexs []*regexp.Regexp
}

func init() {
	processors.Register(New)
}

func New(c interface{}, logger *logrus.Logger) processors.Processor {
	var pc SanatizeConfig
	config.ConfigFromStruct(c, &pc)

	f := &Sanatize{
		logger:          logger,
		trimWhitespaces: pc.TrimWhitespaces,
	}

	for _, regex := range pc.CharacterBlacklistRegexs {
		f.characterReplaceRegexs = append(f.characterReplaceRegexs, regexp.MustCompile(regex))
	}

	return f
}

func (f Sanatize) Process(file *file.File) error {
	file.Filename = f.sanatize(file.Filename)
	return nil
}

func (f Sanatize) sanatize(filename string) string {
	if f.trimWhitespaces {
		filename = strings.TrimSpace(filename)
	}

	for _, removeRegex := range f.characterReplaceRegexs {
		filename = removeRegex.ReplaceAllString(filename, "")
	}
	return filename
}
