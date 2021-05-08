package filter

import (
	"io"

	"github.com/rwese/archivar/archivar/filter/filename"
	"github.com/sirupsen/logrus"
)

type Filter interface {
	Filter(filename, filepath string, data io.Reader) (bool, error)
}

func New(filterType string, config interface{}, logger *logrus.Logger) (filter Filter) {
	switch filterType {
	case "filename":
		filter = filename.New(
			config,
			logger,
		)
	default:
		logger.Panic("could not create new filter from given config")
	}

	return nil
}
