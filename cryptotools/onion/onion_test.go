package onion

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnion(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		onion, err := New()
		assert.NoError(t, err)
		assert.NotNil(t, onion)
		assert.Len(t, onion.PrivKey(), 64)
		assert.Len(t, onion.PubKey(), 32)
	})

	// TODO: test with loading session from file

	t.Run("NewFromPrivKey", func(t *testing.T) {
		// Generate new Onion
		onion, err := New()
		assert.NoError(t, err)

		// Generate new Onion from the private key of the previous one
		onion2, err := NewFromPrivKey(onion.PrivKey())
		assert.NoError(t, err)
		assert.NotNil(t, onion2)

		// Check if the new onion matches the original
		assert.Equal(t, onion.PrivKey(), onion2.PrivKey())
		assert.Equal(t, onion.PubKey(), onion2.PubKey())
		assert.Equal(t, onion.Address(), onion2.Address())
	})

	t.Run("NewFromPrivKey2", func(t *testing.T) {
		keyHex := "f86af341ed3a612ff0754c77b33b60eb1cd40ed204603134217bc857fc411867f2ff4a987ffefcdb3a8c1b5af8f22aff8bcdeb4420426b491747f73fc5327d24"
		expected_onion := "gbislcwjbx2h3pkdavsaqku3mlx4tcnfhmpq2gji5nyegrbrsqcvv6qd"
		keyBytes, err := hex.DecodeString(keyHex)
		assert.NoError(t, err)

		onion, err := NewFromPrivKey(keyBytes)
		assert.NoError(t, err)

		// Check if the new onion matches the original
		assert.Equal(t, onion.Address(), expected_onion)

	})

}
