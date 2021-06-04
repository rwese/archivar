package filesize

import (
	"bytes"
	"errors"
	"io"

	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
)

type FilesizeConfig struct {
	MinSizeBytes int64
	MaxSizeBytes int64
}

type filesize struct {
	MinSizeBytes int64
	MaxSizeBytes int64
	logger       *logrus.Logger
}

func init() {
	filters.Register(New)
}

func New(c interface{}, logger *logrus.Logger) filters.Filter {
	var fc FilesizeConfig
	config.ConfigFromStruct(c, &fc)

	if fc.MaxSizeBytes == 0 && fc.MinSizeBytes == 0 {
		logger.Fatalln("Filesize filter requires at least one MaxSizeBytes/MinSizeBytes")
	}

	f := &filesize{
		logger:       logger,
		MaxSizeBytes: fc.MaxSizeBytes,
		MinSizeBytes: fc.MinSizeBytes,
	}

	return f
}

func (f *filesize) Filter(filterFile *file.File) (result filterResult.Results, err error) {
	var b bytes.Buffer
	var fileSize int64
	sizeReader := io.TeeReader(filterFile.Body, &b)
	if f.MaxSizeBytes > 0 {
		limitReader := io.LimitReader(sizeReader, f.MaxSizeBytes+1)
		readBytes, err := io.ReadAll(limitReader)
		if err != nil {
			return filterResult.Allow, err
		}
		fileSize = int64(len(readBytes))
		if fileSize > f.MaxSizeBytes && !errors.Is(err, io.EOF) {
			f.logger.Debugf("Filesize: Reject MaxSizeBytes %s", filterFile.Filename())
			return filterResult.Reject, nil
		}
	}

	if f.MinSizeBytes >= 0 && fileSize < f.MinSizeBytes {
		leastb := make([]byte, int(f.MinSizeBytes))

		_, err := io.ReadAtLeast(sizeReader, leastb, int(f.MinSizeBytes))
		if err != nil {
			f.logger.Debugf("Filesize: Reject MinSizeBytes %s", filterFile.Filename())
			return filterResult.Reject, nil
		}
	}

	if _, err = io.ReadAll(sizeReader); err != nil {
		return filterResult.Reject, err
	}

	f.logger.Debugf("Filesize: Fallthrough %s", filterFile.Filename())
	filterFile.SetMetadata(file.WithContent(bytes.NewReader(b.Bytes())))

	return filterResult.Allow, nil
}
