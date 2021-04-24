package archivar

import (
	"context"
	"time"

	"github.com/rwese/archivar/archivar/archiver"
	"github.com/rwese/archivar/archivar/gatherer"
	"github.com/sirupsen/logrus"
)

const POLLING_INTERVAL = 60

type Service struct {
	logger *logrus.Logger
	jobs   []ArchiveJob
}

type ArchiveJob struct {
	gatherer gatherer.Gatherer
	archiver archiver.Archiver
}

func New(logger *logrus.Logger) Service {
	return Service{
		logger: logger,
	}
}

func (s *Service) AddJob(gatherer gatherer.Gatherer, archiver archiver.Archiver) {
	s.jobs = append(s.jobs, ArchiveJob{gatherer, archiver})
}

func (s *Service) Run(ctx context.Context) {
	schedule := time.After(1 * POLLING_INTERVAL)

	for {
		select {
		case <-ctx.Done():
			s.logger.Debugln("Gracefully exit")
			s.logger.Debugln(ctx.Err())
			return
		case <-schedule:
			s.logger.Debugln("run...")

			for _, job := range s.jobs {
				err := job.gatherer.Download()
				if err != nil {
					s.logger.Fatalf("error: %s", err.Error())
				}
			}

		}

		schedule = time.After(time.Second * POLLING_INTERVAL)
	}

}
