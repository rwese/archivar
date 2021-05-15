package middleware

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/processor/processors"
	"github.com/rwese/archivar/internal/file"
)

func New(next archivers.Archiver, processor processors.Processor) archivers.Archiver {
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
