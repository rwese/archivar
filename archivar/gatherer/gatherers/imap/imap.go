package imap

import (
	imapClient "github.com/rwese/archivar/internal/imap"
	"github.com/rwese/archivar/utils/config"

	"github.com/emersion/go-imap"
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/gatherer/gatherers"
	"github.com/sirupsen/logrus"
)

type Imap struct {
	DeleteDownloaded bool
	storage          archivers.Archiver
	client           *imapClient.Imap
	logger           *logrus.Logger
}

func init() {
	gatherers.Register(New)
}

func New(c interface{}, storage archivers.Archiver, logger *logrus.Logger) (i gatherers.Gatherer) {
	config.ConfigFromStruct(c, &i)

	return &Imap{
		storage: storage,
		logger:  logger,
		client:  imapClient.New(c, storage, logger),
	}
}

func (i Imap) Download() (err error) {
	done := make(chan error, 1)
	messages := make(chan *imap.Message, 10)
	if err = i.client.GetMessages(messages, done, i.DeleteDownloaded); err != nil {
		return
	}

	if err = <-done; err != nil {
		return
	}

	readMsgSeq := new(imap.SeqSet)
	// readMsgSeq := new(imap.SeqSet)
	for msg := range messages {
		err := i.client.ProcessMessage(*msg, i.storage.Upload)
		if err != nil {
			i.logger.Warnf("Failed to process message: %s", err.Error())
			continue
		}

		readMsgSeq.AddNum(msg.SeqNum)
	}

	if i.DeleteDownloaded {
		i.logger.Debug("deleting processed messages")

		if err = i.client.FlagAndDeleteMessages(readMsgSeq); err != nil {
			i.logger.Fatalf("Failed to clean read messages: %s", err.Error())
		}
	}
	i.logger.Debug("processing imap storage done!")
	return
}
