package imap

import (
	"github.com/rwese/archivar/internal/utils/config"

	"github.com/emersion/go-imap"
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/archiver/archivers/imap/client"
	"github.com/sirupsen/logrus"
)

type ImapGathererConfig struct {
	Server           string
	Username         string
	Password         string
	Inbox            string
	AllowInsecureSSL bool
	DeleteDownloaded bool
}

type ImapGatherer struct {
	deleteDownloaded bool
	storage          archivers.Archiver
	client           *client.Imap
	logger           *logrus.Logger
}

func NewGatherer(c interface{}, storage archivers.Archiver, logger *logrus.Logger) (i archivers.Gatherer) {
	igc := ImapGathererConfig{
		Inbox: "Inbox",
	}

	config.ConfigFromStruct(c, &igc)

	return &ImapGatherer{
		deleteDownloaded: igc.DeleteDownloaded,
		storage:          storage,
		logger:           logger,
		client: client.New(
			igc.Server,
			igc.Username,
			igc.Password,
			igc.Inbox,
			igc.AllowInsecureSSL,
			logger,
		),
	}
}

func (i ImapGatherer) Download() (err error) {
	messages := make(chan *imap.Message, 1)

	readMsgSeq := new(imap.SeqSet)

	go func() {
		for msg := range messages {
			err := i.client.ProcessMessage(msg, i.storage.Upload)
			if err != nil {
				i.logger.Warnf("Failed to process message: %s", err.Error())
				continue
			}

			readMsgSeq.AddNum(msg.SeqNum)
		}
	}()

	if err = i.client.GetMessages(messages, i.deleteDownloaded); err != nil {
		return
	}

	if i.deleteDownloaded {
		i.logger.Debug("deleting processed messages")

		if err = i.client.FlagAndDeleteMessages(readMsgSeq); err != nil {
			i.logger.Fatalf("Failed to clean read messages: %s", err.Error())
		}
	}
	i.logger.Debug("processing imap storage done!")
	return
}

func (i *ImapGatherer) Connect() (err error) {
	if err = i.storage.Connect(); err != nil {
		return
	}

	return i.client.Connect()
}
