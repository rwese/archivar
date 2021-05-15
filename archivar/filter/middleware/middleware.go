package middleware

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters"
	"github.com/rwese/archivar/internal/file"
)

func New(next archivers.Archiver, filter filters.Filter) archivers.Archiver {
	fa := &filterArchiver{next: next, filter: filter}
	return fa
}

type filterArchiver struct {
	archivers.Archiver
	next   archivers.Archiver
	filter filters.Filter
}

func (f *filterArchiver) Upload(file file.File) (err error) {
	result, err := f.filter.Filter(&file)
	if err != nil {
		return err
	}

	if result == filterResult.Reject {
		return nil
	}

	return f.next.Upload(file)
}
