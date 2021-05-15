package filter

import (
	"github.com/rwese/archivar/archivar/filter/filters"
	_ "github.com/rwese/archivar/archivar/filter/filters/filename"
	_ "github.com/rwese/archivar/archivar/filter/filters/filesize"
	"github.com/sirupsen/logrus"
)

// New will return a new filter based on the given typeName and config
func New(filterType string, config interface{}, logger *logrus.Logger) filters.Filter {
	f := filters.Get(filterType, config, logger)

	if f == nil {
		logger.Panicf("could not create new filter '%s' from given config", filterType)
	}

	return f
}
