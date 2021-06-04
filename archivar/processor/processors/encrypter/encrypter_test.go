package encrypter_test

import (
	"bytes"
	"io"
	"testing"

	internalEncrypter "github.com/rwese/archivar/internal/encrypter"

	"github.com/rwese/archivar/archivar/processor/processors/encrypter"
	"github.com/rwese/archivar/internal/file"
	"github.com/sirupsen/logrus"
)

// These are just testing keys do not worry for my safety
var publicKey = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAMT8Us5GBUj9DFpJf0JCeCkx+dhOFV/T0YqCn827yY//iMOhgD+LEeR7
BKYLeJ24G06QDc0DotEw1gjHSYqNl7x2qnDJfKoOvZzdgh0SERfMUju1K7kxSM+W
Wnq7gv6vwyG4VMkWZapVtW/vVGiJFCXW8oLHGgdyiG35szvPatGDAgMBAAE=
-----END RSA PUBLIC KEY-----
`

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDE/FLORgVI/QxaSX9CQngpMfnYThVf09GKgp/Nu8mP/4jDoYA/
ixHkewSmC3iduBtOkA3NA6LRMNYIx0mKjZe8dqpwyXyqDr2c3YIdEhEXzFI7tSu5
MUjPllp6u4L+r8MhuFTJFmWqVbVv71RoiRQl1vKCxxoHcoht+bM7z2rRgwIDAQAB
AoGANan/7Q4KVo4JlXc8YhK1pZNl21W6YPbVuQRJAMVN7hrRaWpQA/+hCjuxUoMB
gwYq+kYoXFfHPXIufQm9sS9NzKGxnTVAXeA+AWhM259i/mtbXw0fWS41FCNjWUEn
LG9pk/D1ccuf6fNo+e/RIa1gCAnhJyi3dIR2vSePzJMnV7ECQQDdk/zH+dZ77Wod
sLM/4Ut094xGImCsJBEAk+AasCUn5SzUMWyO8OAuwawSlJxG4kGL/gW6n5bXhCw7
xJAcz0XNAkEA45ZPqOBa1ZqTwmUomx4xeuqyEKA9dJCCAy+jiiJKgohsGA1ppdup
CjRa6+7HNEMT5eFJ5AjrHXNUF0a8ZxokjwJBAIzKPXIrc3dnEWgwIJVUaAe4S288
5MQ8XnlJfLo4dkN1QRjLFrl0oF3VParItsvrc86p56X/RW9HUnvfl9pWcXkCQQCr
bIjMF1HUGv65KiEP1gpHH4jIZSplJoQHilaQsYuWDtP8uf2d5HrLKOxjUhPSFcRj
HvLdRKp0IG5yqeE3d8WZAkBGv9Q31J6l5L2WxNzRG9uxkdAKZcvKzxq4+NRak8tk
/+TAnM6quG+EvJNPx8PQ/ll78PqiCqJKzuSO5LgvrbRd
-----END RSA PRIVATE KEY-----`

func TestEncryptTrim(t *testing.T) {
	fileTests := map[string]struct {
		config encrypter.EncrypterConfig
		have   file.File
		want   file.File
	}{
		"encrypt": {
			config: encrypter.EncrypterConfig{
				AddExtension: ".encrypted",
				PublicKey:    publicKey,
			},
			have: *file.New(
				file.WithContent(bytes.NewReader([]byte(` Testing `))),
				file.WithFilename("a1b2"),
			),
			want: *file.New(
				file.WithContent(bytes.NewReader([]byte(` Testing `))),
				file.WithFilename("a1b2.encrypted"),
			),
		},
	}

	for testName, fileTest := range fileTests {
		f := encrypter.New(fileTest.config, logrus.New())

		file := fileTest.have
		f.Process(&file)
		if file.Filename() != fileTest.want.Filename() {
			t.Fatalf("Failed test '%s'", testName)
		}

		haveContent, _ := io.ReadAll(file.Body)
		pkey, _ := internalEncrypter.DecodePrivateKey([]byte(privateKey))
		e := internalEncrypter.New(nil, pkey)
		decryptedContent, _ := e.Decrypt(haveContent)

		wantContent, _ := io.ReadAll(fileTest.want.Body)

		if !bytes.Equal(wantContent, decryptedContent) {
			t.Fatalf("Encryption - Decryption failed '%s'", testName)
		}
	}
}
