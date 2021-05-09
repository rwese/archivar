package archivar

import (
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/filter"
	"github.com/rwese/archivar/archivar/gatherer"
	"github.com/rwese/archivar/archivar/job"
	"github.com/sirupsen/logrus"
)

type Archivar struct {
	logger *logrus.Logger
	jobs   []job.Job
	config Config
}

type ConfigSub struct {
	Interval int
	Type     string
	Config   interface{}
}

type Config struct {
	Settings struct {
		DefaultInterval int
		Log             struct {
			Debugging bool
		}
	}
	Archivers map[string]ConfigSub
	Gatherers map[string]ConfigSub
	Filters   map[string]ConfigSub
	Jobs      map[string]job.JobsConfig
}

func New(config Config, logger *logrus.Logger) Archivar {
	s := Archivar{
		logger: logger,
		config: config,
	}

	for jobName, job := range config.Jobs {
		s.addJob(jobName, job)
	}

	return s
}

func (s *Archivar) addJob(jobName string, job job.JobsConfig) {
	interval := s.config.Settings.DefaultInterval
	if job.Interval != 0 {
		interval = job.Interval
	}

	c := s.config.Archivers[job.Archiver]
	archiver := archiver.New(c.Type, c.Config, s.logger)

	for _, filterName := range job.Filters {
		c = s.config.Filters[filterName]
		f := filter.New(c.Type, c.Config, s.logger)
		archiver = filter.FilterArchiverMiddleware(archiver, f)
	}

	c = s.config.Gatherers[job.Gatherer]
	gatherer := gatherer.New(c.Type, c.Config, archiver, s.logger)
	s.AddJob(jobName, interval, gatherer)
}
