package job

import (
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/gatherer"
	"github.com/sirupsen/logrus"
)

type Job struct {
	Name     string
	Interval int
	Errors   int
	Gatherer gatherer.Gatherer
	Archiver archiver.Archiver
	Logger   *logrus.Logger
}

func (j *Job) Download() error {
	return j.Gatherer.Download()
}

type JobsConfig struct {
	Interval int
	Gatherer string
	Archiver string
	Filters  []string
}
