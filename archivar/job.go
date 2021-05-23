package archivar

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/filter/filters"
	"github.com/rwese/archivar/archivar/job"
	"github.com/rwese/archivar/archivar/processor/processors"
)

func (s *Archivar) runJob(job job.Job, ctx context.Context, stop context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	s.logger.Infof("Starting Job %s, every %d seconds", job.Name, job.Interval)

	waitingTime := time.Duration(job.Interval) * time.Second
	schedule := time.After(time.Second * 0)

	if err := job.Connect(); err != nil {
		s.logger.Warnf("%s: error %s", job.Name, err.Error())
		stop()
	}

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

func (s *Archivar) RunJobs() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := new(sync.WaitGroup)
	wg.Add(len(s.jobs))

	for _, job := range s.jobs {
		go s.runJob(job, ctx, stop, wg)
	}

	wg.Wait()
}

func (s *Archivar) Dump() {
	fmt.Println("Archivers:")
	for _, a := range archivers.ListArchivers() {
		fmt.Printf("\t%s\n", a)
	}
	fmt.Println("Gatherers:")
	for _, a := range archivers.ListGatherers() {
		fmt.Printf("\t%s\n", a)
	}
	fmt.Println("Filters:")
	for _, a := range filters.ListFilters() {
		fmt.Printf("\t%s\n", a)
	}
	fmt.Println("Processors:")
	for _, a := range processors.ListProcessors() {
		fmt.Printf("\t%s\n", a)
	}
}
