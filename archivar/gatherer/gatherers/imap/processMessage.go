package imap

import (
	"fmt"
	"io"
	"log"
	"path"
	"regexp"
	"time"

	"github.com/emersion/go-imap"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/sirupsen/logrus"
)

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

func (i Imap) processMessage(msg *imap.Message) (err error) {
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

			cleanSubject := subjectCleanup.ReplaceAllString(mailData.subject, "")
			if len(cleanSubject) > SUBJECT_LENGTH {
				cleanSubject = cleanSubject[:SUBJECT_LENGTH]
			}
			filename := cleanSubject + fileExt
			if err = i.storage.Upload(filename, filePrefixPath, p.Body); err != nil {
				return err
			}
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			if filename == "" {
				continue
			}

			logrus.Debugf("Got attachment: %v", filename)
			logrus.Debugf("Saving as: %v", filePrefixPath)

			if err = i.storage.Upload(filename, filePrefixPath, p.Body); err != nil {
				return err
			}
		}
	}
	return
}
