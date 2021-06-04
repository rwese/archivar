package processor

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/processor/processors"
	"github.com/rwese/archivar/internal/file"
)

func NewMiddleware(next archivers.Archiver, processor processors.Processor) archivers.Archiver {
	fa := &processorArchiver{next: next, processor: processor}
	return fa
}

type processorArchiver struct {
	archivers.Archiver
	next      archivers.Archiver
	processor processors.Processor
}

func (f *processorArchiver) Upload(file *file.File) (err error) {
	err = f.processor.Process(file)
	if err != nil {
		return err
	}

	return f.next.Upload(file)
}

func (f *processorArchiver) Connect() (err error) {
	return f.next.Connect()
}
