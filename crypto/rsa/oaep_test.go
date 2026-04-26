package rsa

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"strings"
	"testing"
)

func generateRSAKeyPairDER(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key failed: %v", err)
	}

	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key failed: %v", err)
	}

	privateDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal private key failed: %v", err)
	}

	return publicDER, privateDER
}

func generateECDSAKeyPairDER(t *testing.T) ([]byte, []byte) {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate ecdsa key failed: %v", err)
	}

	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal ecdsa public key failed: %v", err)
	}

	privateDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal ecdsa private key failed: %v", err)
	}

	return publicDER, privateDER
}

func TestEncryptOAEPDecryptOAEP_Success(t *testing.T) {
	publicDER, privateDER := generateRSAKeyPairDER(t)
	plainText := []byte("hello toolbox rsa oaep")
	label := []byte("toolbox")

	cipherText, err := EncryptOAEP(plainText, publicDER, sha256.New(), label)
	if err != nil {
		t.Fatalf("EncryptOAEP failed: %v", err)
	}
	if len(cipherText) == 0 {
		t.Fatal("EncryptOAEP returned empty ciphertext")
	}

	decrypted, err := DecryptOAEP(cipherText, privateDER, sha256.New(), label)
	if err != nil {
		t.Fatalf("DecryptOAEP failed: %v", err)
	}

	if !bytes.Equal(decrypted, plainText) {
		t.Fatalf("decrypted text mismatch: got %q want %q", decrypted, plainText)
	}
}

func TestEncryptOAEPBase64DecryptOAEPBase64_Success(t *testing.T) {
	publicDER, privateDER := generateRSAKeyPairDER(t)
	plainText := []byte("base64 round trip")

	cipherTextBase64, err := EncryptOAEPBase64(plainText, publicDER, sha512.New(), nil)
	if err != nil {
		t.Fatalf("EncryptOAEPBase64 failed: %v", err)
	}
	if strings.TrimSpace(cipherTextBase64) == "" {
		t.Fatal("EncryptOAEPBase64 returned empty string")
	}

	decrypted, err := DecryptOAEPBase64(cipherTextBase64, privateDER, sha512.New(), nil)
	if err != nil {
		t.Fatalf("DecryptOAEPBase64 failed: %v", err)
	}

	if !bytes.Equal(decrypted, plainText) {
		t.Fatalf("decrypted text mismatch: got %q want %q", decrypted, plainText)
	}
}

func TestEncryptOAEP_InvalidPublicKey(t *testing.T) {
	_, privateDER := generateRSAKeyPairDER(t)

	if _, err := EncryptOAEP([]byte("x"), []byte("invalid"), sha256.New(), nil); err == nil {
		t.Fatal("EncryptOAEP should fail for invalid public key der")
	}

	ecdsaPublicDER, _ := generateECDSAKeyPairDER(t)
	_, err := EncryptOAEP([]byte("x"), ecdsaPublicDER, sha256.New(), nil)
	if err == nil || err.Error() != "invalid public key" {
		t.Fatalf("EncryptOAEP should return invalid public key error, got: %v", err)
	}

	if _, err := DecryptOAEP([]byte("x"), privateDER, sha256.New(), nil); err == nil {
		t.Fatal("DecryptOAEP should fail for malformed ciphertext")
	}
}

func TestDecryptOAEP_InvalidPrivateKey(t *testing.T) {
	publicDER, _ := generateRSAKeyPairDER(t)
	cipherText, err := EncryptOAEP([]byte("hello"), publicDER, sha256.New(), nil)
	if err != nil {
		t.Fatalf("EncryptOAEP failed: %v", err)
	}

	if _, err := DecryptOAEP(cipherText, []byte("invalid"), sha256.New(), nil); err == nil {
		t.Fatal("DecryptOAEP should fail for invalid private key der")
	}

	_, ecdsaPrivateDER := generateECDSAKeyPairDER(t)
	_, err = DecryptOAEP(cipherText, ecdsaPrivateDER, sha256.New(), nil)
	if err == nil || err.Error() != "invalid private key" {
		t.Fatalf("DecryptOAEP should return invalid private key error, got: %v", err)
	}
}

func TestDecryptOAEP_LabelMismatch(t *testing.T) {
	publicDER, privateDER := generateRSAKeyPairDER(t)
	cipherText, err := EncryptOAEP([]byte("hello"), publicDER, sha256.New(), []byte("label-a"))
	if err != nil {
		t.Fatalf("EncryptOAEP failed: %v", err)
	}

	if _, err := DecryptOAEP(cipherText, privateDER, sha256.New(), []byte("label-b")); err == nil {
		t.Fatal("DecryptOAEP should fail when label does not match")
	}
}

func TestDecryptOAEP_HashMismatch(t *testing.T) {
	publicDER, privateDER := generateRSAKeyPairDER(t)
	cipherText, err := EncryptOAEP([]byte("hello"), publicDER, sha256.New(), nil)
	if err != nil {
		t.Fatalf("EncryptOAEP failed: %v", err)
	}

	if _, err := DecryptOAEP(cipherText, privateDER, sha512.New(), nil); err == nil {
		t.Fatal("DecryptOAEP should fail when hash does not match")
	}
}

func TestEncryptOAEP_PlainTextLengthBoundary(t *testing.T) {
	publicDER, privateDER := generateRSAKeyPairDER(t)

	// For RSA-2048 with SHA-256, max OAEP message size is 256 - 2*32 - 2 = 190 bytes.
	maxPlainText := bytes.Repeat([]byte{'a'}, 190)
	tooLongPlainText := bytes.Repeat([]byte{'b'}, 191)

	cipherText, err := EncryptOAEP(maxPlainText, publicDER, sha256.New(), nil)
	if err != nil {
		t.Fatalf("EncryptOAEP should succeed on max boundary, got: %v", err)
	}

	decrypted, err := DecryptOAEP(cipherText, privateDER, sha256.New(), nil)
	if err != nil {
		t.Fatalf("DecryptOAEP failed: %v", err)
	}
	if !bytes.Equal(decrypted, maxPlainText) {
		t.Fatal("decrypted max-boundary text mismatch")
	}

	if _, err := EncryptOAEP(tooLongPlainText, publicDER, sha256.New(), nil); err == nil {
		t.Fatal("EncryptOAEP should fail for plaintext longer than OAEP limit")
	}
}

func TestDecryptOAEPBase64_InvalidInput(t *testing.T) {
	_, privateDER := generateRSAKeyPairDER(t)
	if _, err := DecryptOAEPBase64("@@not-base64@@", privateDER, sha256.New(), nil); err == nil {
		t.Fatal("DecryptOAEPBase64 should fail for invalid base64 input")
	}
}

