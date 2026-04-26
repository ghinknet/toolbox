package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"hash"
)

func EncryptOAEP(plainText []byte, block []byte, hash hash.Hash, label []byte) ([]byte, error) {
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}

	cipherText, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, plainText, label)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func EncryptOAEPBase64(plainText []byte, block []byte, hash hash.Hash, label []byte) (string, error) {
	data, err := EncryptOAEP(plainText, block, hash, label)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecryptOAEP(cipherText []byte, block []byte, hash hash.Hash, label []byte) ([]byte, error) {
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block)
	if err != nil {
		return nil, err
	}

	privateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key")
	}

	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, cipherText, label)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func DecryptOAEPBase64(cipherText string, block []byte, hash hash.Hash, label []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	return DecryptOAEP(data, block, hash, label)
}
