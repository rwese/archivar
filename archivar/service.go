package archivar

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/filter"
	"github.com/rwese/archivar/archivar/gatherer"
	"github.com/rwese/archivar/archivar/job"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	jobs   []job.Job
}

type ConfigSub struct {
	Type   string
	Config interface{}
}

type Config struct {
	Settings struct {
		Log struct {
			Debugging bool
		}
	}
	Archivers map[string]*ConfigSub
	Gatherers map[string]*ConfigSub
	Filters   map[string]*ConfigSub
	Jobs      map[string]ConfigJobs
}

type ConfigJobs struct {
	Gatherer string
	Archiver string
	Filters  []string
}

func New(config Config, logger *logrus.Logger) Service {
	s := Service{
		logger: logger,
	}

	for _, job := range config.Jobs {

		var filters []filter.Filter
		var c *ConfigSub
		for _, filterName := range job.Filters {
			c = config.Filters[filterName]
			f := filter.New(c.Type, c.Config, s.logger)
			filters = append(filters, f)
		}

		c = config.Archivers[job.Archiver]
		archiver := archiver.New(c.Type, c.Config, s.logger)

		c = config.Gatherers[job.Gatherer]
		gatherer := gatherer.New(c.Type, c.Config, archiver, filters, s.logger)
		s.AddJob(gatherer, archiver)
	}

	return s
}

func (s *Service) AddJob(gatherer gatherer.Gatherer, archiver archiver.Archiver) {
	s.jobs = append(s.jobs, job.Job{
		Gatherer: gatherer,
		Archiver: archiver,
	})
}

func (s *Service) run() {
	wg := new(sync.WaitGroup)
	wg.Add(len(s.jobs))

	for _, job := range s.jobs {
		go s.runJob(job, wg)
	}

	wg.Wait()
}

func (s *Service) runJob(job job.Job, wg *sync.WaitGroup) {
	defer wg.Done()
	err := job.Download()
	if err != nil {
		s.logger.Fatalf("error: %s", err.Error())
	}
}

func (s *Service) Watch(pollingInterval int) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	waitingTime := time.Duration(pollingInterval) * time.Second
	schedule := time.After(time.Second * 0)

	for {
		select {
		case <-ctx.Done():
			s.logger.Debugln("Gracefully exit")
			return
		case <-schedule:
			s.run()
		}

		if waitingTime <= 0 {
			s.logger.Debugln("Stopping after single run")
			break
		}

		time.Sleep(time.Second * 1)

		// time.After allows the process to be stopped instantly which sleep doesn't
		schedule = time.After(waitingTime)
	}
}
