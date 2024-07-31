package client

import (
	"bytes"
	"crypto/tls"
	id "github.com/emersion/go-imap-id"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"
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
	inboxPrefix      string
	processingInbox  string
	allowInsecureSSL bool
	timestampFormat  string
	pathPattern      string
	filePattern      string
	maxSubjectLength int64
	withDeleted      bool
	withSeen         bool
	client           *client.Client
	section          *imap.BodySectionName
	items            []imap.FetchItem
	logger           *logrus.Logger
}

func New(server, username, password, inbox, inboxPrefix string, allowInsecureSSL bool, timestampFormat string, pathPattern string, filePattern string, maxSubjectLength int64, withSeen bool, withDeleted bool, logger *logrus.Logger) *Imap {
	i := &Imap{
		server:           server,
		username:         username,
		password:         password,
		inbox:            inbox,
		inboxPrefix:      inboxPrefix,
		allowInsecureSSL: allowInsecureSSL,
		timestampFormat:  timestampFormat,
		pathPattern:      pathPattern,
		filePattern:      filePattern,
		maxSubjectLength: maxSubjectLength,
		withSeen:         withSeen,
		withDeleted:      withDeleted,
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

	// set id information
	idClient := id.NewClient(i.client)
	_, err = idClient.ID(
		id.ID{id.FieldName: "IMAPClient", id.FieldVersion: "1.2.0"}, // just define it casually and declare your identity
	)
	if err != nil {
		i.logger.Warnf("failed to set client ID information, %s", err.Error())
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

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		log.Fatal(err)
	}

	readerFirst := bytes.NewReader(buf.Bytes())
	readerSecond := bytes.NewReader(buf.Bytes())

	m, err := mail.CreateReader(readerFirst)
	if err != nil {
		log.Fatal(err)
	}
	mailData := mailData{}
	header := m.Header

	if date, err := header.Date(); err == nil {
		mailData.date = date
	}
	if from, err := header.AddressList("From"); err == nil {
		mailData.from = NewParsedAddress(from[0])
	}
	if to, err := header.AddressList("To"); err == nil {
		mailData.to = NewParsedAddress(to[0])
	}
	if subject, err := header.Subject(); err == nil {
		if i.maxSubjectLength > 0 && int64(len(subject)) > i.maxSubjectLength {
			subject = subject[:i.maxSubjectLength]
		}

		mailData.subject = subject
	}

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

		filePath := mailData.getFilePath(i.processingInbox, i.pathPattern, i.timestampFormat)
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

			filename = fileSafe.ReplaceAllString(mailData.subject, "") + fileExt

		case *mail.AttachmentHeader:
			filename, _ = h.Filename()
			if filename == "" {
				continue
			}

			logrus.Debugf("Got attachment: %v", filename)
			logrus.Debugf("Saving as: %v", filePath)

		default:
			continue
		}

		filename = mailData.getFileName(filename, i.filePattern, i.timestampFormat)

		f := file.New(
			file.WithContent(p.Body),
			file.WithFilename(fileSafe.ReplaceAllString(filename, "")),
			file.WithDirectory(filePath),
			file.WithCreatedAt(mailData.date),
		)

		if err = upload(f); err != nil {
			return err
		}
	}

	// save eml file
	emlFileName := fileSafe.ReplaceAllString(mailData.subject, "") + ".eml"
	emlFilePath := mailData.getFilePath(i.processingInbox, i.pathPattern, i.timestampFormat)
	emlFile := file.New(
		file.WithContent(readerSecond),
		file.WithFilename(fileSafe.ReplaceAllString(emlFileName, "")),
		file.WithDirectory(emlFilePath),
		file.WithCreatedAt(mailData.date),
	)
	if err = upload(emlFile); err != nil {
		return err
	}

	return nil
}

type mailData struct {
	date    time.Time
	from    ParsedAddress
	to      ParsedAddress
	subject string
}

type namingPatternMapping struct {
	Placeholder string
	Value       string
}

type namingPattern struct {
	mapping []namingPatternMapping
}

type ParsedAddress struct {
	Address *mail.Address
	Full    string
	User    string
	Detail  string
	Domain  string
}

func NewParsedAddress(a *mail.Address) ParsedAddress {
	at := strings.LastIndex(a.Address, "@")
	domain := ""
	if at > 0 {
		domain = a.Address[at+1:]
	}

	user := a.Address
	detail := ""
	mailToSubMatch := emailPlusPart.FindSubmatch([]byte(a.String()))
	if len(mailToSubMatch) > 1 {
		user = strings.ReplaceAll(a.Address, string(mailToSubMatch[0]), "@")
		detail = string(mailToSubMatch[1])
	} else if at > 0 {
		user = a.Address[:at]
	}

	user = strings.ReplaceAll(user, "@"+domain, "")

	return ParsedAddress{
		Address: a,
		Full:    a.Address,
		User:    user,
		Detail:  detail,
		Domain:  domain,
	}
}

func (n *namingPattern) format(format string) string {
	s := format
	for _, p := range n.mapping {
		search := "{" + p.Placeholder + "}"
		s = strings.Replace(s, search, p.Value, -1)
	}

	return s
}

func (n *namingPattern) add(placeholder string, value string) {
	n.mapping = append(n.mapping, namingPatternMapping{
		Placeholder: placeholder,
		Value:       value,
	})
}

var emailPlusPart = regexp.MustCompile(`\+(.+?)\@`)
var fileSafe = regexp.MustCompile(`[\\/:"*?<>|]+`)

func (m mailData) getNamingPattern(dateFormat string) namingPattern {
	patternList := namingPattern{}
	patternList.add("mail_from", m.from.Full)
	patternList.add("mail_from_user", m.from.User)
	patternList.add("mail_from_detail", m.from.Detail)
	patternList.add("mail_from_domain", m.from.Domain)
	patternList.add("mail_to", m.to.Full)
	patternList.add("mail_to_user", m.to.User)
	patternList.add("mail_to_detail", m.to.Detail)
	patternList.add("mail_to_domain", m.to.Domain)
	patternList.add("mail_subject", m.subject)
	patternList.add("mail_subject_safe", fileSafe.ReplaceAllString(m.subject, ""))
	patternList.add("mail_date", m.date.Format(dateFormat))

	return patternList
}

func (m mailData) getFilePath(inbox string, pathPattern string, timestampFormat string) string {
	pattern := m.getNamingPattern(timestampFormat)
	pattern.add("mail_dir", inbox)

	p := pattern.format(pathPattern)
	p = filepath.Clean(p)

	return p
}

func (m mailData) getFileName(fileName string, fileNamePattern string, timestampFormat string) string {
	pattern := m.getNamingPattern(timestampFormat)
	pattern.add("attachment_filename", fileName)

	return pattern.format(fileNamePattern)
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

func (i *Imap) MoveMessages(seqSet *imap.SeqSet, dest string) (err error) {
	return i.client.Move(seqSet, dest)
}

func (i *Imap) getInboxesByPrefix(prefix string) []string {
	i.Connect()
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- i.client.List("", prefix+"*", mailboxes)
	}()

	inboxes := []string{}
	for m := range mailboxes {
		inboxes = append(inboxes, m.Name)
	}

	return inboxes
}
func (i *Imap) ListInboxes() {
	i.Connect()
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- i.client.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}
}

func (i *Imap) GetMessages(messageChan chan *imap.Message) (err error) {
	i.Connect()

	inboxes := []string{}
	if i.inbox != "" {
		inboxes = append(inboxes, i.inbox)
	}

	if i.inboxPrefix != "" {
		inboxes = append(inboxes, i.getInboxesByPrefix(i.inboxPrefix)...)
	}

	for _, inbox := range inboxes {
		i.processingInbox = inbox

		err := i.processInboxMessages(inbox, messageChan)
		if err != nil {
			i.logger.Warnf("failed to process inbox message %v", err)
			return nil
		}
	}

	return nil
}

func (i *Imap) processInboxMessages(inbox string, messageChan chan *imap.Message) (err error) {
	_, err = i.client.Select(inbox, false)
	if err != nil {
		i.ListInboxes()
		close(messageChan)
		return err
	}

	i.logger.Debugf("selected '%s'", inbox)

	criteria := imap.NewSearchCriteria()
	if i.withSeen == false {
		criteria.WithoutFlags = append(criteria.WithoutFlags, imap.SeenFlag)
	}
	if i.withDeleted == false {
		criteria.WithoutFlags = append(criteria.WithoutFlags, imap.DeletedFlag)
	}

	foundMsgs, err := i.client.Search(criteria)
	if err != nil {
		close(messageChan)
		return err
	}

	if len(foundMsgs) == 0 {
		i.logger.Debug("no messages")
		close(messageChan)
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
