package filesize

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/rwese/archivar/archivar/filter/filterResult"
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

func New(config interface{}, logger *logrus.Logger) *Filesize {
	jsonM, _ := json.Marshal(config)
	var c FilesizeConfig
	json.Unmarshal(jsonM, &c)

	if c.MaxSizeBytes == 0 && c.MinSizeBytes == 0 {
		logger.Fatalln("Filesize filter requires at least one MaxSizeBytes/MinSizeBytes")
	}

	f := &Filesize{
		logger:       logger,
		MaxSizeBytes: c.MaxSizeBytes,
		MinSizeBytes: c.MinSizeBytes,
	}

	return f
}

func (f *Filesize) Filter(filename *string, filepath *string, data *io.Reader) (filterResult.Results, error) {
	r, err := io.ReadAll(*data)
	if err != nil {
		return filterResult.Allow, err
	}

	*data = bytes.NewReader(r)
	fileSize := len(r)
	if f.MaxSizeBytes > 0 && fileSize > f.MaxSizeBytes {
		f.logger.Debugf("Reject Filesize-MaxSizeBytes %s", *filename)
		return filterResult.Reject, nil
	}

	if f.MinSizeBytes > 0 && fileSize < f.MinSizeBytes {
		f.logger.Debugf("Reject Filesize-MinSizeBytes %s", *filename)
		return filterResult.Reject, nil
	}

	f.logger.Debugf("Fallthrough Filesize %s", *filename)
	return filterResult.Allow, err
}
