package auth

import (
	"bytes"
	"testing"

	myaes "github.com/1F47E/go-shaihulud/cryptotools/aes"

	"github.com/stretchr/testify/assert"
)

// func TestPasswordGen(t *testing.T) {
// 	// example password
// 	// 3688-7BE9
// 	// TODO: test password gen
// }

func TestAuthEncryptDecrypt(t *testing.T) {
	t.Run("Encryption and Decryption Gen pass", func(t *testing.T) {
		// Create a new AES crypter
		crypter := &myaes.AEScrypter{}

		auth, err := New(crypter, "")
		assert.NoError(t, err, "Error creating new Auth instance")

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

func TestOnionFromKey(t *testing.T) {
	key := "AF3EAFDE09FBA80741034641180F13E029B056BC5F7440598EAC2EBFFE894D6C51D5263782D957FC95A856E1469159BFC97228448D2BF5F2DC896CE25758EF742235A7CEA5032C3F0B0B8A78EB8B08BA7D036E436F563078E660ED46"
	password := "F2A6-D23A"
	expected_onion := "nachzaurfkn742gnigmm6aqkjqubmojwykvcuenzt53423e775vcinid"
	crypter := myaes.New()
	auth, err := NewFromKey(crypter, key, password)
	assert.NoError(t, err, "Error creating new Auth instance")
	assert.Equal(t, auth.OnionAddress(), expected_onion)
}
