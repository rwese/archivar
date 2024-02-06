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
	InboxPrefix      string
	AllowInsecureSSL bool
	DeleteDownloaded bool
	TimestampFormat  string
	PathPattern      string
	FilePattern      string
	MaxSubjectLength int64
}

type ImapGatherer struct {
	deleteDownloaded bool
	storage          archivers.Archiver
	client           *client.Imap
	logger           *logrus.Logger
}

func NewGatherer(c interface{}, storage archivers.Archiver, logger *logrus.Logger) (i archivers.Gatherer) {
	igc := ImapGathererConfig{}

	config.ConfigFromStruct(c, &igc)

	if igc.Inbox == "" && igc.InboxPrefix == "" {
		igc.Inbox = "Inbox"
	}

	if igc.TimestampFormat == "" {
		igc.TimestampFormat = "20060102_150405"
	}

	if igc.PathPattern == "" {
		igc.PathPattern = "{mail_dir}/{mail_to}/{mail_to_detail}/{mail_date}-{mail_subject_safe}"
	}

	if igc.FilePattern == "" {
		igc.FilePattern = "{attachment_filename}"
	}

	return &ImapGatherer{
		deleteDownloaded: igc.DeleteDownloaded,
		storage:          storage,
		logger:           logger,
		client: client.New(
			igc.Server,
			igc.Username,
			igc.Password,
			igc.Inbox,
			igc.InboxPrefix,
			igc.AllowInsecureSSL,
			igc.TimestampFormat,
			igc.PathPattern,
			igc.FilePattern,
			igc.MaxSubjectLength,
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

	if i.deleteDownloaded && !readMsgSeq.Empty() {
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
