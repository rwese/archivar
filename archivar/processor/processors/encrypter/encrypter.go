package encrypter

import (
	"bytes"
	"crypto/rsa"
	"io"

	internalEncrypter "github.com/rwese/archivar/internal/encrypter"

	"github.com/rwese/archivar/archivar/processor/processors"
	"github.com/rwese/archivar/internal/file"
	"github.com/rwese/archivar/internal/utils/config"
	"github.com/sirupsen/logrus"
)

type EncrypterConfig struct {
	AddExtension string
	DontRename   bool
	PublicKey    string
}

type Encrypter struct {
	encryptedFileExtension string
	dontRename             bool
	key                    *rsa.PublicKey
	logger                 *logrus.Logger
	client                 internalEncrypter.Encrypter
}

func init() {
	processors.Register(New)
}

func New(c interface{}, logger *logrus.Logger) processors.Processor {
	var pc EncrypterConfig
	pc.AddExtension = ".encrypted"

	config.ConfigFromStruct(c, &pc)

	publicKey, err := internalEncrypter.DecodePublicKey([]byte(pc.PublicKey))
	if err != nil {
		logger.Fatal(err)
	}

	f := &Encrypter{
		logger:                 logger,
		client:                 internalEncrypter.New(publicKey, nil),
		key:                    publicKey,
		encryptedFileExtension: pc.AddExtension,
		dontRename:             pc.DontRename,
	}

	return f
}

func (f Encrypter) Process(file *file.File) (err error) {
	encrypted, err := f.encrypt(&file.Body)
	if err != nil {
		return
	}

	file.Body = bytes.NewReader(encrypted)

	if !f.dontRename {
		file.Filename = file.Filename + f.encryptedFileExtension
	}

	return nil
}

func (f Encrypter) encrypt(body *io.Reader) (encrypted []byte, err error) {
	bodyData, err := io.ReadAll(*body)
	if err != nil {
		return
	}

	return f.client.Encrypt(bodyData)
}
