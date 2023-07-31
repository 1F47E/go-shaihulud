package myrsa

import (
	"bytes"
	"testing"

	msgcrypter "go-dmtor/cryptotools/message_crypter"
)

func TestRSACrypter(t *testing.T) {
	rsa, err := New()
	if err != nil {
		t.Fatalf("rsa init error: %v\n", err)
	}
	crypter := msgcrypter.New(rsa)
	pubKeyBytes := crypter.PubKey()

	message := []byte("Hello World!")

	// Encrypt a message
	cipher, err := crypter.Encrypt(message, pubKeyBytes)
	if err != nil {
		t.Fatalf("encrypt error: %v\n", err)
	}

	// Decrypt the message
	plain, err := crypter.Decrypt(cipher)
	if err != nil {
		t.Fatalf("decrypt error: %v\n", err)
	}

	// test assert orig and decoded
	if !bytes.Equal(plain, message) {
		t.Fatalf("orig and decoded do not match: %s != %s\n", plain, message)
	}
}
