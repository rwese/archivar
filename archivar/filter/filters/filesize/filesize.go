package filesize

import (
	"bytes"
	"errors"
	"io"

	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/archivar/filter/filters"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/utils/config"
	"github.com/sirupsen/logrus"
)

type FilesizeConfig struct {
	MinSizeBytes int64
	MaxSizeBytes int64
}

type Filesize struct {
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

	f := &Filesize{
		logger:       logger,
		MaxSizeBytes: fc.MaxSizeBytes,
		MinSizeBytes: fc.MinSizeBytes,
	}

	return f
}

func (f *Filesize) Filter(file *file.File) (result filterResult.Results, err error) {
	var buffer bytes.Buffer
	var fileSize int64
	sizeReader := io.MultiWriter(&buffer)
	if f.MaxSizeBytes > 0 {
		fileSize, err = io.CopyN(sizeReader, file.Body, f.MaxSizeBytes+1)
		if fileSize >= f.MaxSizeBytes && !errors.Is(err, io.EOF) {
			f.logger.Debugf("Filesize: Reject MaxSizeBytes %s", file.Filename)
			return filterResult.Reject, nil
		}
	} else {
		fileSize, err = io.Copy(sizeReader, file.Body)
	}

	if f.MinSizeBytes >= 0 && fileSize < f.MinSizeBytes {
		f.logger.Debugf("Filesize: Reject MinSizeBytes %s", file.Filename)
		result = filterResult.Reject
	} else {
		f.logger.Debugf("Filesize: Fallthrough %s", file.Filename)
		result = filterResult.Allow
		file.Body = bytes.NewReader(buffer.Bytes())
	}

	return
}
