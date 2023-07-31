package auth

import (
	"bytes"
	myaes "go-dmtor/cryptotools/aes"
	"go-dmtor/cryptotools/onion"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthEncryptDecrypt(t *testing.T) {
	t.Run("Encryption and Decryption", func(t *testing.T) {
		// Create a new AES crypter
		crypter := &myaes.AEScrypter{}

		// Create onion instance
		oni, err := onion.New()
		assert.NoError(t, err)

		// Create a new Auth instance with empty onion
		auth := New(crypter, oni)

		// The plaintext we want to encrypt
		plaintext := []byte("This is some test plaintext")

		// Attempt to encrypt the plaintext
		ciphertext, err := auth.Encrypt(plaintext)
		assert.NoError(t, err, "Error encrypting plaintext")

		// Attempt to decrypt the ciphertext
		decrypted, err := auth.Decrypt(ciphertext)
		assert.NoError(t, err, "Error decrypting ciphertext")

		// Check if the decrypted text matches the original plaintext
		assert.True(t, bytes.Equal(decrypted, plaintext), "Decrypted text does not match original plaintext")
	})
}
