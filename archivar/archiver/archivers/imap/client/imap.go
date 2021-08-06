package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"path"
	"regexp"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/rwese/archivar/archivar/archiver/archivers"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

type Imap struct {
	server           string
	username         string
	password         string
	inbox            string
	allowInsecureSSL bool
	client           *client.Client
	section          *imap.BodySectionName
	items            []imap.FetchItem
	logger           *logrus.Logger
}

func New(server, username, password, inbox string, allowInsecureSSL bool, logger *logrus.Logger) *Imap {
	i := &Imap{
		server:           server,
		username:         username,
		password:         password,
		inbox:            inbox,
		allowInsecureSSL: allowInsecureSSL,
		logger:           logger,
	}

	i.section = &imap.BodySectionName{}
	i.items = []imap.FetchItem{i.section.FetchItem()}
	imap.CharsetReader = charset.Reader

	return i
}

func (i *Imap) Connect() (err error) {
	if i.client != nil {
		state := i.client.State()
		if state == imap.AuthenticatedState || state == imap.SelectedState {
			return
		}
	}

	tlsConfig := tls.Config{InsecureSkipVerify: i.allowInsecureSSL}
	i.logger.Debugf("connecting to %s", i.server)
	i.client, err = client.DialTLS(i.server, &tlsConfig)
	if err != nil {
		i.logger.Fatalf("failed to connect to imap: %s", err.Error())
	}

	i.logger.Debugf("authenticate as %s using password %t", i.username, i.password != "")
	if err = i.client.Login(i.username, i.password); err != nil {
		i.logger.Fatalf("failed to login to imap: %s", err.Error())
	}

	return
}
func (i *Imap) Disconnect() (err error) {
	return i.client.Logout()
}

func (i Imap) ProcessMessage(msg *imap.Message, upload archivers.UploadFunc) error {
	r := msg.GetBody(i.section)

	m, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}
	mailData := mailData{}
	header := m.Header

	if date, err := header.Date(); err == nil {
		mailData.date = date
	}
	if from, err := header.AddressList("From"); err == nil {
		mailData.from = *from[0]
	}
	if to, err := header.AddressList("To"); err == nil {
		mailData.to = *to[0]
	}
	if subject, err := header.Subject(); err == nil {
		mailData.subject = subject
	}

	filePrefixPath := mailData.getFilePath()

	if err != nil {
		log.Fatal(err)
	}

	for {
		p, err := m.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		var filename string
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()

			fileExt := ""
			if contentType == "text/html" {
				fileExt = ".html"
			} else if contentType == "text/plain" {
				fileExt = ".txt"
			} else {
				logrus.Debugf("Skipping Content-Type: %s", contentType)
				continue
			}

			filename = mailData.subject + fileExt

		case *mail.AttachmentHeader:
			filename, _ = h.Filename()
			if filename == "" {
				continue
			}

			logrus.Debugf("Got attachment: %v", filename)
			logrus.Debugf("Saving as: %v", filePrefixPath)

		default:
			continue
		}

		f := file.New(
			file.WithContent(p.Body),
			file.WithFilename(filename),
			file.WithDirectory(filePrefixPath),
			file.WithCreatedAt(mailData.date),
		)

		if err = upload(f); err != nil {
			return err
		}
	}
	return nil
}

type mailData struct {
	date    time.Time
	from    mail.Address
	to      mail.Address
	subject string
}

var emailPlusPart = regexp.MustCompile(`\+(.+?)\@`)
var subjectCleanup = regexp.MustCompile(`[^a-zA-Z0-9\-_ ]+`)

const SUBJECT_LENGTH = 30

func (m mailData) getFilePath() string {
	// TODO add variant options
	timestamp := fmt.Sprintf(
		"%04d%02d%02d_%02d%02d%02d",
		m.date.Year(),
		m.date.Month(),
		m.date.Day(),
		m.date.Hour(),
		m.date.Minute(),
		m.date.Second(),
	)

	rootDirectory := m.to.Address
	foundPlusString := emailPlusPart.FindSubmatch([]byte(m.to.String()))
	if len(foundPlusString) > 1 {
		rootDirectory = string(foundPlusString[1])
	}

	subjectCleanupPath := subjectCleanup.ReplaceAllString(m.subject, "")

	return path.Join(rootDirectory, timestamp+"_"+subjectCleanupPath)
}

func (i Imap) FlagAndDeleteMessages(readMsgSeq *imap.SeqSet) (err error) {
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err := i.client.Store(readMsgSeq, item, flags, nil); err != nil {
		return err
	}

	if err := i.client.Expunge(nil); err != nil {
		return err
	}

	return
}

func (i *Imap) GetMessages(messageChan chan *imap.Message, deleteDownloaded bool) (err error) {
	i.Connect()

	mbox, err := i.client.Select(i.inbox, false)
	if err != nil {
		return
	}

	i.logger.Debugf("selected '%s'", i.inbox)

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.DeletedFlag}
	foundMsgs, err := i.client.Search(criteria)
	if mbox.Messages == 0 {
		i.logger.Debug("no messages")
		return nil
	}

	i.logger.Debugf("found %d messages", len(foundMsgs))
	seqset := new(imap.SeqSet)
	seqset.AddNum(foundMsgs...)

	go func() {
		if err := i.client.Fetch(seqset, i.items, messageChan); err != nil {
			i.logger.Warnf("failed to fetch from imap %v", err)
		}
	}()

	return nil
}
