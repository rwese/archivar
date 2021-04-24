package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rwese/archivar/archivar"
	"github.com/rwese/archivar/archivar/archiver/webdav"
	"github.com/rwese/archivar/archivar/gatherer/imap"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logDebugging := os.Getenv("LOG_DEBUGGING")
	logger.SetLevel(logrus.InfoLevel)

	if logDebugging == "1" {
		logger.SetLevel(logrus.DebugLevel)
	}

	webdavUploader := webdav.New(
		os.Getenv("WEBDAV_URL"),
		os.Getenv("WEBDAV_USERNAME"),
		os.Getenv("WEBDAV_PASSWORD"),
		os.Getenv("WEBDAV_UPLOAD_DIRECTORY"),
		logger,
	)

	imapDownloader := imap.New(
		os.Getenv("IMAP_SERVER"),
		os.Getenv("IMAP_USERNAME"),
		os.Getenv("IMAP_PASSWORD"),
		os.Getenv("SETTINGS_KEEPUPLOADED") == "1",
		webdavUploader,
		logger,
	)

	s := archivar.New(logger)
	s.AddJob(imapDownloader, webdavUploader)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	s.Run(ctx)
}
