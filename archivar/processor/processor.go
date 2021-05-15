package processor

import (
	"github.com/rwese/archivar/archivar/processor/processors"
	_ "github.com/rwese/archivar/archivar/processor/processors/sanatizer"
	"github.com/sirupsen/logrus"
)

func New(processorType string, config interface{}, logger *logrus.Logger) processors.Processor {
	p := processors.Get(processorType, config, logger)

	if p == nil {
		logger.Panicf("could not create new processor '%s' from given config", processorType)
	}

	return p
}
