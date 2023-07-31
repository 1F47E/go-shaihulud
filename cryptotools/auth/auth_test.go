package auth

import (
	"bytes"
	myaes "go-dmtor/cryptotools/aes"
	"testing"
)

func TestAuthEncryptDecrypt(t *testing.T) {
	// Create a new AES crypter
	crypter := &myaes.AEScrypter{}

	// Create a new Auth instance
	a := New(crypter)

	// The plaintext we want to encrypt
	plaintext := []byte("This is some test plaintext")

	// The password we're using to encrypt the plaintext
	password := "ThisIsASecurePassword"

	// Attempt to encrypt the plaintext
	ciphertext, err := a.Encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Error encrypting plaintext: %v", err)
	}

	// Attempt to decrypt the ciphertext
	decrypted, err := a.Decrypt(ciphertext, password)
	if err != nil {
		t.Fatalf("Error decrypting ciphertext: %v", err)
	}

	// Check if the decrypted text matches the original plaintext
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypted text does not match original plaintext")
	}
}
