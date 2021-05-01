package imap

import (
	"crypto/tls"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/rwese/archivar/archivar/archiver"
	"github.com/sirupsen/logrus"
)

type Imap struct {
	storage      archiver.Archiver
	server       string
	username     string
	password     string
	keepUploaded bool
	inbox        string
	section      *imap.BodySectionName
	items        []imap.FetchItem
	logger       *logrus.Logger
}

func New(server string, username string, password string, keepUploaded bool, storage archiver.Archiver, logger *logrus.Logger) (i *Imap) {
	i = &Imap{
		storage:      storage,
		server:       server,
		username:     username,
		password:     password,
		keepUploaded: keepUploaded,
		inbox:        "INBOX",
		logger:       logger,
	}

	i.section = &imap.BodySectionName{}
	i.items = []imap.FetchItem{i.section.FetchItem()}
	return i
}

func (i *Imap) Connect() (c *client.Client, err error) {
	tlsConfig := tls.Config{InsecureSkipVerify: true}
	i.logger.Debugf("connecting to %s", i.server)
	c, err = client.DialTLS(i.server, &tlsConfig)
	if err != nil {
		i.logger.Fatalf("failed to connect to imap: %s", err.Error())
	}

	i.logger.Debugf("authenticate as %s using password %t", i.username, i.password != "")
	if err = c.Login(i.username, i.password); err != nil {
		i.logger.Fatalf("failed to login to imap: %s", err.Error())
	}

	return
}

func (i *Imap) Download() (err error) {
	c, err := i.Connect()
	defer c.Logout()
	mbox, err := c.Select(i.inbox, false)
	if err != nil {
		i.logger.Fatal(err)
	}
	i.logger.Debugf("selected '%s'", i.inbox)

	if mbox.Messages == 0 {
		i.logger.Debug("no messages")
		return nil
	}
	i.logger.Debugf("found %d messages", mbox.Messages)
	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, i.items, messages)
	}()

	readMsgSeq := new(imap.SeqSet)
	for msg := range messages {
		err = i.processMessage(msg)

		if err != nil {
			i.logger.Warnf("Failed to process message: %s", err.Error())
		} else {
			readMsgSeq.AddNum(msg.SeqNum)
		}
	}

	if err := <-done; err != nil {
		i.logger.Fatal(err)
		return err
	}

	if !i.keepUploaded {
		i.logger.Debug("deleting processed messages")

		if err = i.flagAndDeleteMessages(readMsgSeq, c); err != nil {
			i.logger.Fatalf("Failed to clean read messages: %s", err.Error())
		}
	}

	i.logger.Debug("processing imap storage done!")
	return
}

func (i Imap) flagAndDeleteMessages(readMsgSeq *imap.SeqSet, c *client.Client) (err error) {
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err := c.Store(readMsgSeq, item, flags, nil); err != nil {
		return err
	}

	if err := c.Expunge(nil); err != nil {
		return err
	}

	return
}
