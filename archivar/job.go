package archivar

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rwese/archivar/archivar/gatherer/gatherers"
	"github.com/rwese/archivar/archivar/job"
)

func (s *Archivar) AddJob(jobName string, interval int, gatherer gatherers.Gatherer) {
	s.jobs = append(s.jobs, job.Job{
		Name:     jobName,
		Interval: interval,
		Gatherer: gatherer,
	})
}

func (s *Archivar) runJob(job job.Job, ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	s.logger.Infof("Starting Job %s, every %d seconds", job.Name, job.Interval)

	waitingTime := time.Duration(job.Interval) * time.Second
	schedule := time.After(time.Second * 0)

	for {
		select {
		case <-ctx.Done():
			s.logger.Debugf("%s: Gracefully exit", job.Name)
			return
		case <-schedule:
			s.logger.Debugf("%s: Run job", job.Name)
			err := job.Download()
			if err != nil {
				s.logger.Warnf("%s: error %s", job.Name, err.Error())
				stop()
			}
		}

		if waitingTime <= 0 {
			s.logger.Debugf("%s: Stopping after single run", job.Name)
			break
		}

		time.Sleep(time.Second * 1)

		// time.After allows the process to be stopped instantly which sleep doesn't
		schedule = time.After(waitingTime)
	}

	s.logger.Debugf("%s: ended", job.Name)
}

func (s *Archivar) RunJobs(pollingInterval int) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := new(sync.WaitGroup)
	wg.Add(len(s.jobs))

	for _, job := range s.jobs {
		go s.runJob(job, ctx, stop, wg)
	}

	wg.Wait()
}
