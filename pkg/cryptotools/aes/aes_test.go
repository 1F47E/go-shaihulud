package caes

import (
	"bytes"
	"testing"

	"go-dmtor/pkg/cryptotools/auth"

	"github.com/stretchr/testify/assert"
)

func TestAESCrypter(t *testing.T) {
	t.Run("Encryption and Decryption", func(t *testing.T) {
		var myaes = new(AEScrypter)

		a, err := auth.New(myaes, "")
		assert.NoError(t, err)

		// TODO: test with loading session from file

		data := []byte("Hello World!")

		// Encrypt a message
		cipher, err := a.Encrypt(data)
		assert.NoError(t, err, "Encrypt error")

		// Decrypt the message
		plain, err := a.Decrypt(cipher)
		assert.NoError(t, err, "Decrypt error")

		// Test assert original and decoded
		assert.True(t, bytes.Equal(plain, data), "Original and decoded do not match: %s != %s", plain, data)
	})
}
