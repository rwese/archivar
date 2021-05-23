package job

import (
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/sirupsen/logrus"
)

type Job struct {
	Name     string
	Interval int
	Errors   int
	Gatherer archivers.Gatherer
	Logger   *logrus.Logger
}

// Download fill perform the Gatherer.Download() which will use the Archiver
func (j *Job) Download() error {
	return j.Gatherer.Download()
}

// Connect will connect the Gatherer and Archiver to verify the configuration
func (j *Job) Connect() (err error) {
	if err = j.Gatherer.Connect(); err != nil {
		return
	}

	return
}

type JobConfig struct {
	Interval   int
	Gatherer   string
	Archiver   string
	Filters    []string
	Processors []string
}
