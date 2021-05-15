package processor

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/processor/processors"
	_ "github.com/rwese/archivar/archivar/processor/processors/sanatizer"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

func New(processorType string, config interface{}, logger *logrus.Logger) processors.Processor {
	p := processors.Get(processorType, config, logger)

	if p == nil {
		logger.Panicf("could not create new processor '%s' from given config", processorType)
	}

	return p
}

func ProcessorArchiverMiddleware(next archivers.Archiver, processor processors.Processor) archivers.Archiver {
	fa := &ProcessorArchiver{next: next, processor: processor}
	return fa
}

type ProcessorArchiver struct {
	archivers.Archiver
	next      archivers.Archiver
	processor processors.Processor
}

func (f *ProcessorArchiver) Upload(file file.File) (err error) {
	err = f.processor.Process(&file)
	if err != nil {
		return err
	}

	return f.next.Upload(file)
}
