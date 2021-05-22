package encrypter

// Encrypter is my implementation of using PublicKey Crypto in tandem with AES
// Problems I had were that rsa.EncryptPKCS1v15 max length was limited by the
// KeyLength, and AES required a passphrase as it is a symetric cipher.
//
// My idea is it to generate a random password and encrypt the password for the
// given publicKey so the password can be decrypted only with the privateKey.
//
// References where I got "inspritation" and stolen code from:
//
// *  @stupidbodo https://gist.github.com/stupidbodo/601b68bfef3449d1b8d9
//

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rwese/archivar/internal/random"
)

type Encrypter struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Encrypter {
	return Encrypter{publicKey: publicKey, privateKey: privateKey}
}

func (e Encrypter) SplitFile(src string) (err error) {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	passphraseEncrypted := input[:128]
	data := input[128:]

	dstKeyFile := src + ".rsa.key"
	dstDataFile := src + ".aes.body"
	if err = os.WriteFile(dstKeyFile, passphraseEncrypted, 0660); err != nil {
		return
	}

	fmt.Printf("write encrypted key to: %s\n", dstKeyFile)

	err = os.WriteFile(dstDataFile, data, 0660)

	fmt.Printf("write encrypted data to: %s\n", dstDataFile)

	return
}

func (e Encrypter) DecryptFile(src, dst string) (err error) {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	output, err := e.Decrypt(input)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, output, 0660)
}

func (e Encrypter) EncryptFile(dst, src string) (err error) {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	output, err := e.Encrypt(input)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, output, 0660)
}

func (e Encrypter) Encrypt(data []byte) ([]byte, error) {
	encryptionKey := random.RandomBytes(32)
	encryptedEncryptionKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, e.publicKey, []byte(encryptionKey), nil)
	if err != nil {
		return nil, err
	}

	encryptedData, err := encrypt(encryptionKey, string(data))
	if err != nil {
		return nil, err
	}
	encryptedEncryptionKey = []byte(removeBase64Padding(base64.URLEncoding.EncodeToString(encryptedEncryptionKey)))
	joinedOutput := bytes.Join(
		[][]byte{
			encryptedEncryptionKey,
			[]byte(encryptedData),
		},
		[]byte(""),
	)

	return joinedOutput, nil
}

func (e Encrypter) Decrypt(data []byte) ([]byte, error) {
	passphraseEncrypted := data[:171]
	passphraseEncrypted, err := base64.URLEncoding.DecodeString(addBase64Padding(string(passphraseEncrypted)))
	if err != nil {
		return nil, err
	}

	data = data[171:]
	passphrase, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, e.privateKey, passphraseEncrypted, nil)
	if err != nil {
		return nil, err
	}

	decrypted, err := decrypt(passphrase, string(data))
	if err != nil {
		return nil, err
	}

	return []byte(decrypted), nil
}

func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func encrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	msg := Pad([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(msg))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(msg))
	finalMsg := removeBase64Padding(base64.URLEncoding.EncodeToString(ciphertext))
	return finalMsg, nil
}

func decrypt(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(addBase64Padding(text))
	if err != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	iv := decodedMsg[:aes.BlockSize]
	msg := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}
