package sanatizer

import (
	"encoding/json"
	"io"
	"regexp"
	"strings"

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

func New(config interface{}, logger *logrus.Logger) *Sanatize {
	jsonM, _ := json.Marshal(config)
	var fc SanatizeConfig
	json.Unmarshal(jsonM, &fc)

	f := &Sanatize{
		logger:          logger,
		trimWhitespaces: fc.TrimWhitespaces,
	}

	for _, regex := range fc.CharacterBlacklistRegexs {
		f.characterReplaceRegexs = append(f.characterReplaceRegexs, regexp.MustCompile(regex))
	}

	return f
}

func (f *Sanatize) Process(filename *string, filepath *string, data *io.Reader) error {
	*filename = f.sanatize(*filename)
	return nil
}

func (f *Sanatize) sanatize(filename string) string {
	if f.trimWhitespaces {
		filename = strings.TrimSpace(filename)
	}

	for _, removeRegex := range f.characterReplaceRegexs {
		filename = removeRegex.ReplaceAllString(filename, "")
	}
	return filename
}
