package myrsa

import (
	"testing"

	msgcrypter "github.com/1F47E/go-shaihulud/internal/cryptotools/msgcrypter"

	"github.com/stretchr/testify/require"
)

func TestRSACrypter(t *testing.T) {
	rsa, err := New()
	require.NoError(t, err, "RSA initialization failed")

	crypter := msgcrypter.New(rsa)
	pubKeyBytes := crypter.PubKey()

	t.Run("Test Message Encryption and Decryption", func(t *testing.T) {
		message := []byte("Hello World!")

		// Encrypt a message
		cipher, err := crypter.Encrypt(message, pubKeyBytes)
		require.NoError(t, err, "Encryption failed")

		// Decrypt the message
		plain, err := crypter.Decrypt(cipher)
		require.NoError(t, err, "Decryption failed")

		// Test assert orig and decoded
		require.Equal(t, message, plain, "Original and decoded message do not match")
	})

}
