package caes

import (
	"bytes"
	"testing"

	"go-dmtor/cryptotools/auth"
	"go-dmtor/cryptotools/onion"

	"github.com/stretchr/testify/assert"
)

func TestAESCrypter(t *testing.T) {
	t.Run("Encryption and Decryption", func(t *testing.T) {
		var myaes = new(AEScrypter)

		oni, err := onion.New()
		assert.NoError(t, err)

		crypter := auth.New(myaes, oni)

		data := []byte("Hello World!")

		// Encrypt a message
		cipher, err := crypter.Encrypt(data)
		assert.NoError(t, err, "Encrypt error")

		// Decrypt the message
		plain, err := crypter.Decrypt(cipher)
		assert.NoError(t, err, "Decrypt error")

		// Test assert original and decoded
		assert.True(t, bytes.Equal(plain, data), "Original and decoded do not match: %s != %s", plain, data)
	})
}
