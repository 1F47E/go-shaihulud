package caes

import (
	"bytes"
	"testing"

	"go-dmtor/cryptotools/auth"
)

func TestAESCrypter(t *testing.T) {
	var myaes = new(AEScrypter)
	crypter := auth.New(myaes)
	password := "myPassword"

	data := []byte("Hello World!")

	// Encrypt a message
	cipher, err := crypter.Encrypt(data, password)
	if err != nil {
		t.Fatalf("encrypt error: %v\n", err)
	}

	// Decrypt the message
	plain, err := crypter.Decrypt(cipher, password)
	if err != nil {
		t.Fatalf("decrypt error: %v\n", err)
	}

	// Test assert original and decoded
	if !bytes.Equal(plain, data) {
		t.Fatalf("original and decoded do not match: %s != %s\n", plain, data)
	}
}
