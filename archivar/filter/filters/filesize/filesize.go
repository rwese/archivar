package filesize

import (
	"bytes"
	"io"

	"github.com/rwese/archivar/archivar/filter/filterResult"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/utils/config"
	"github.com/sirupsen/logrus"
)

type FilesizeConfig struct {
	MinSizeBytes int
	MaxSizeBytes int
}

type Filesize struct {
	MinSizeBytes int
	MaxSizeBytes int
	logger       *logrus.Logger
}

func New(c interface{}, logger *logrus.Logger) *Filesize {
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
	if f.MaxSizeBytes > 0 {
		file.Body = io.LimitReader(file.Body, int64(f.MaxSizeBytes)+1)
	}

	r, err := io.ReadAll(file.Body)
	if err != nil {
		return filterResult.Allow, err
	}

	fileSize := len(r)
	file.Body = bytes.NewReader(r)

	if f.MaxSizeBytes > 0 && fileSize > f.MaxSizeBytes {
		f.logger.Debugf("Filesize: Reject MaxSizeBytes %s", file.Filename)
		result = filterResult.Reject
	} else if f.MinSizeBytes > 0 && fileSize < f.MinSizeBytes {
		f.logger.Debugf("Filesize: Reject MinSizeBytes %s", file.Filename)
		result = filterResult.Reject
	} else {
		f.logger.Debugf("Filesize: Fallthrough %s", file.Filename)
		result = filterResult.Allow
	}

	return
}
