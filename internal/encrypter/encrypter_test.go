package encrypter_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/rwese/archivar/internal/encrypter"
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

func TestEncrypt(t *testing.T) {
	testCases := map[string]struct {
		have []byte
		want []byte
	}{
		"encrypt": {
			have: []byte("a"),
			want: []byte("DvymHk4CJQ6g/NDwM9dfulxa8vhgNEBMnudREK3EIezuI2aFM9" +
				"R/TaViqB/u92q1a9J/6X4NZXQjUKX+TArVG6ar4ALDPrmdiY5lF0TgjTulzskNB" +
				"WRNm23JdoG+9uGawYjQhONpvgK9j2G5CbU+4w3J9Y9TqWj8Kw+oJYU4hYpKTWdJ" +
				"Mk43ei0xRWRRWUZwa1I3X25GclVOcmdHNTYwckxJSlFTN3B1Z1lV"),
		},
	}

	for _, testCase := range testCases {
		pbkey, _ := encrypter.DecodePublicKey([]byte(publicKey))
		pkey, _ := encrypter.DecodePrivateKey([]byte(privateKey))
		e := encrypter.New(pbkey, pkey)
		haveEnc, err := e.Encrypt(testCase.have)
		if err != nil {
			t.Fatal(err)
		}
		haveEncBase64 := base64.StdEncoding.EncodeToString([]byte(haveEnc))
		haveDenc, err := e.Decrypt(haveEnc)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(testCase.have, haveDenc) {
			t.Fatal("Encryption does not match Decryption")
		}

		if bytes.Equal(testCase.want, []byte(haveEncBase64)) {
			t.Fatal("Encryption missmatch")
		}
	}
}
