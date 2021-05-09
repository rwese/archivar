package filter

import (
	"io"

	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters/filename"
	"github.com/rwese/archivar/archivar/filter/filters/sanatize"
	"github.com/sirupsen/logrus"
)

type Filter interface {
	Filter(filename *string, filepath *string, data io.Reader) (filterResult.Results, error)
}

func New(filterType string, config interface{}, logger *logrus.Logger) (filter Filter) {
	switch filterType {
	case "filename":
		return filename.New(
			config,
			logger,
		)
	case "sanatize":
		return sanatize.New(
			config,
			logger,
		)
	default:
		logger.Panic("could not create new filter from given config")
	}

	return nil
}

func FilterArchiverMiddleware(next archiver.Archiver, filter Filter) archiver.Archiver {
	fa := &FilterArchiver{next: next, filter: filter}
	return fa
}

type FilterArchiver struct {
	archiver.Archiver
	next   archiver.Archiver
	filter Filter
}

func (f *FilterArchiver) Upload(fileName string, fileDirectory string, fileHandle io.Reader) (err error) {
	result, err := f.filter.Filter(&fileName, &fileDirectory, fileHandle)
	if err != nil {
		return err
	}

	if result == filterResult.Reject {
		return nil
	}

	return f.next.Upload(fileName, fileDirectory, fileHandle)
}
