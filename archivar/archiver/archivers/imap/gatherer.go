package imap

import (
	"github.com/emersion/go-imap"
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/archivar/archiver/archivers/imap/client"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
	"sync"
)

type ImapGathererConfig struct {
	Server                string
	Username              string
	Password              string
	Inbox                 string
	InboxPrefix           string
	AllowInsecureSSL      bool
	DeleteDownloaded      bool
	MoveProcessedToFolder string
	TimestampFormat       string
	PathPattern           string
	FilePattern           string
	MaxSubjectLength      int64
	WithSeen              bool
	WithDeleted           bool
}

type ImapGatherer struct {
	deleteDownloaded      bool
	moveProcessedToFolder string
	storage               archivers.Archiver
	client                *client.Imap
	logger                *logrus.Logger
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
		deleteDownloaded:      igc.DeleteDownloaded,
		moveProcessedToFolder: igc.MoveProcessedToFolder,
		storage:               storage,
		logger:                logger,
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
			igc.WithSeen,
			igc.WithDeleted,
			logger,
		),
	}
}

func (i ImapGatherer) Download() (err error) {
	messages := make(chan *imap.Message, 1)

	readMsgSeq := new(imap.SeqSet)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for msg := range messages {
			err := i.client.ProcessMessage(msg, i.storage.Upload)
			if err != nil {
				i.logger.Warnf("Failed to process message: %s", err.Error())
				continue
			}

			readMsgSeq.AddNum(msg.SeqNum)
		}
	}()

	if err = i.client.GetMessages(messages); err != nil {
		return
	}

	wg.Wait()

	if !readMsgSeq.Empty() {

		if i.deleteDownloaded {
			i.logger.Debug("deleting processed messages")

			if err = i.client.FlagAndDeleteMessages(readMsgSeq); err != nil {
				i.logger.Fatalf("Failed to clean read messages: %s", err.Error())
			}
		} else if len(i.moveProcessedToFolder) > 0 {
			dest := i.moveProcessedToFolder
			err = i.client.MoveMessages(readMsgSeq, dest)
			if err != nil {
				i.logger.Fatalf("Failed to move processed messages to '%s': %s", dest, err.Error())
				return
			}

			i.logger.Debug("Moved mail messages to: " + dest)
		}
	}

	err = i.client.Disconnect()
	if err != nil {
		i.logger.Debugf("Failed to disconnect IMAP %s", err.Error())
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
