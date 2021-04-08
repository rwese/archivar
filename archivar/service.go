package archivar

import (
	"context"
	"os"
	"time"

	"github.com/rwese/archivar/archivar/archiver/webdav"
	"github.com/rwese/archivar/archivar/gatherer/imap"
	"github.com/sirupsen/logrus"
)

const POLLING_INTERVAL = 60

type Service struct {
	logger *logrus.Logger
}

func New() Service {
	logger := logrus.New()

	logDebugging := os.Getenv("LOG.DEBUGGING")
	logger.SetLevel(logrus.InfoLevel)

	if logDebugging == "1" {
		logger.SetLevel(logrus.DebugLevel)
	}

	return Service{
		logger: logger,
	}
}
func (s *Service) Run(ctx context.Context) {
	webdavUploader := webdav.New(os.Getenv("WEBDAV.URL"),
		os.Getenv("WEBDAV.USERNAME"),
		os.Getenv("WEBDAV.PASSWORD"),
		os.Getenv("WEBDAV.UPLOAD.DIRECTORY"),
		s.logger,
	)

	imapDownloader := imap.New(
		os.Getenv("IMAP.SERVER"),
		os.Getenv("IMAP.USERNAME"),
		os.Getenv("IMAP.PASSWORD"),
		os.Getenv("SETTINGS.KEEPUPLOADED") == "1",
		webdavUploader,
		s.logger,
	)

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("Gracefully exit")
			s.logger.Println(ctx.Err())
			return
		default:
			s.logger.Debugln("run...")

			err := imapDownloader.Download()
			if err != nil {
				s.logger.Fatalf("error: %s", err.Error())
			}

			time.Sleep(time.Second * POLLING_INTERVAL)
		}
	}

}
