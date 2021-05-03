package archivar

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/gatherer"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger          *logrus.Logger
	pollingInterval time.Duration
	jobs            []ArchivarJob
}

type Config struct {
	Settings struct {
		Log struct {
			Debugging bool
		}
	}
	Archivers map[string]archiver.ArchiverConfig
	Gatherers map[string]gatherer.GathererConfig
	Jobs      map[string]Jobs
}

type Jobs struct {
	Gatherer string
	Archiver string
}

type ArchivarJob struct {
	gatherer gatherer.Gatherer
	archiver archiver.Archiver
}

func New(config Config, logger *logrus.Logger) Service {
	s := Service{
		logger: logger,
	}

	for _, job := range config.Jobs {
		archiverConfig := config.Archivers[job.Archiver]
		archiver, err := archiver.New(archiverConfig, s.logger)
		if err != nil {
			s.logger.Fatalln(err)
		}

		gathererConfig := config.Gatherers[job.Gatherer]
		gatherer, err := gatherer.New(gathererConfig, archiver, s.logger)
		if err != nil {
			s.logger.Fatalln(err)
		}

		s.AddJob(gatherer, archiver)
	}

	return s
}

func (s *Service) AddJob(gatherer gatherer.Gatherer, archiver archiver.Archiver) {
	s.jobs = append(s.jobs, ArchivarJob{gatherer, archiver})
}

func (s *Service) run() {
	wg := new(sync.WaitGroup)
	wg.Add(len(s.jobs))

	for _, job := range s.jobs {
		go s.runJob(job, wg)
	}

	wg.Wait()
}

func (s *Service) runJob(job ArchivarJob, wg *sync.WaitGroup) {
	defer wg.Done()
	err := job.gatherer.Download()
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

		time.Sleep(time.Second * 1)

		// time.After allows the process to be stopped instantly which sleep doesn't
		schedule = time.After(waitingTime)
	}
}
