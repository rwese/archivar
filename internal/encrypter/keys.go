package encrypter

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateKey(keysize int) (err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keysize)
	if err != nil {
		return
	}

	publicKey := privateKey.PublicKey
	PrintPEMPrivateKey(privateKey)
	PrintPEMPublicKey(&publicKey)

	SavePEMPublicKey("public.pem", &publicKey)
	SavePEMPrivateKey("private.pem", privateKey)

	return
}

func SavePEMPublicKey(fileName string, key *rsa.PublicKey) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer outFile.Close()
	pemString, err := EncodePublicKey(key)
	if err != nil {
		return
	}

	_, err = outFile.Write([]byte(pemString))
	return
}

func SavePEMPrivateKey(fileName string, key *rsa.PrivateKey) (err error) {
	outFile, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer outFile.Close()
	pemString, err := EncodePrivateKey(key)
	if err != nil {
		return
	}

	_, err = outFile.Write([]byte(pemString))
	return
}

func PrintPEMPrivateKey(key *rsa.PrivateKey) (err error) {
	pemString, err := EncodePrivateKey(key)
	if err != nil {
		return
	}
	fmt.Print(string(pemString))
	return
}

func PrintPEMPublicKey(key *rsa.PublicKey) (err error) {
	pemString, err := EncodePublicKey(key)
	if err != nil {
		return
	}
	fmt.Print(string(pemString))
	return
}

func LoadPublicKeyFile(fileName string) (publicKey *rsa.PublicKey, err error) {
	keyFile, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	return DecodePublicKey(keyFile)
}

func LoadPrivateKeyFile(fileName string) (privateKey *rsa.PrivateKey, err error) {
	keyFile, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	return DecodePrivateKey(keyFile)
}

func DecodePublicKey(key []byte) (publicKey *rsa.PublicKey, err error) {
	decoder, _ := pem.Decode(key)
	return x509.ParsePKCS1PublicKey(decoder.Bytes)
}

func DecodePrivateKey(key []byte) (privateKey *rsa.PrivateKey, err error) {
	decoder, _ := pem.Decode(key)
	return x509.ParsePKCS1PrivateKey(decoder.Bytes)
}

func EncodePublicKey(key *rsa.PublicKey) ([]byte, error) {
	var publicKey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(key),
	}
	var outputKey bytes.Buffer
	err := pem.Encode(&outputKey, publicKey)
	return outputKey.Bytes(), err
}

func EncodePrivateKey(key *rsa.PrivateKey) ([]byte, error) {
	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	var outputKey bytes.Buffer
	err := pem.Encode(&outputKey, privateKey)
	return outputKey.Bytes(), err
}
