package filter

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters/filename"
	"github.com/rwese/archivar/archivar/filter/filters/filesize"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

type Filter interface {
	Filter(*file.File) (filterResult.Results, error)
}

func New(filterType string, config interface{}, logger *logrus.Logger) (filter Filter) {
	switch filterType {
	case "filename":
		return filename.New(
			config,
			logger,
		)
	case "filesize":
		return filesize.New(
			config,
			logger,
		)
	default:
		logger.Panicf("could not create new filter '%s' from given config", filterType)
	}

	return nil
}

func FilterArchiverMiddleware(next archivers.Archiver, filter Filter) archivers.Archiver {
	fa := &FilterArchiver{next: next, filter: filter}
	return fa
}

type FilterArchiver struct {
	archivers.Archiver
	next   archivers.Archiver
	filter Filter
}

func (f *FilterArchiver) Upload(file file.File) (err error) {
	result, err := f.filter.Filter(&file)
	if err != nil {
		return err
	}

	if result == filterResult.Reject {
		return nil
	}

	return f.next.Upload(file)
}
