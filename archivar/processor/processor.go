package processor

import (
	"io"

	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/processor/processors/sanatizer"
	"github.com/sirupsen/logrus"
)

type Processor interface {
	Process(filename *string, filepath *string, data *io.Reader) error
}

func New(processorType string, config interface{}, logger *logrus.Logger) (processor Processor) {
	switch processorType {
	case "sanatize":
		return sanatizer.New(
			config,
			logger,
		)
	default:
		logger.Panicf("could not create new processor '%s' from given config", processorType)
	}

	return nil
}

func ProcessorArchiverMiddleware(next archiver.Archiver, processor Processor) archiver.Archiver {
	fa := &ProcessorArchiver{next: next, processor: processor}
	return fa
}

type ProcessorArchiver struct {
	archiver.Archiver
	next      archiver.Archiver
	processor Processor
}

func (f *ProcessorArchiver) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	err = f.processor.Process(&fileName, &fileDirectory, &fileHandle)
	if err != nil {
		return err
	}

	return f.next.Upload(fileName, fileDirectory, fileHandle)
}
