package job

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/gatherer/gatherers"
	"github.com/sirupsen/logrus"
)

type Job struct {
	Name     string
	Interval int
	Errors   int
	Gatherer gatherers.Gatherer
	Archiver archivers.Archiver
	Logger   *logrus.Logger
}

func (j *Job) Download() error {
	return j.Gatherer.Download()
}

type JobsConfig struct {
	Interval   int
	Gatherer   string
	Archiver   string
	Filters    []string
	Processors []string
}
