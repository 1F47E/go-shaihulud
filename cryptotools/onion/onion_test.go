package onion

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnion(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		onion, err := New()
		assert.NoError(t, err)
		assert.NotNil(t, onion)
		assert.Len(t, onion.PrivateKey(), 64)
		assert.Len(t, onion.PubKey(), 32)
	})

	t.Run("NewFromSession", func(t *testing.T) {
		// Generate new Onion and save to file
		onion, err := New()
		assert.NoError(t, err)
		err = onion.Save()
		assert.NoError(t, err)

		// Try to load the onion from the saved session file
		onion2, err := NewFromSession(onion.Address())
		assert.NoError(t, err)
		assert.NotNil(t, onion2)

		// Check if the loaded onion matches the original
		assert.Equal(t, onion.PrivateKey(), onion2.PrivateKey())
		assert.Equal(t, onion.PubKey(), onion2.PubKey())
		assert.Equal(t, onion.Address(), onion2.Address())

		// Clean up session file
		err = os.Remove(SESSION_DIR + "/" + onion.Session())
		assert.NoError(t, err)

		// remove all files
		err = os.RemoveAll(SESSION_DIR)
		assert.NoError(t, err)

	})

	t.Run("NewFromPrivKey", func(t *testing.T) {
		// Generate new Onion
		onion, err := New()
		assert.NoError(t, err)

		// Generate new Onion from the private key of the previous one
		onion2, err := NewFromPrivKey(onion.PrivateKey())
		assert.NoError(t, err)
		assert.NotNil(t, onion2)

		// Check if the new onion matches the original
		assert.Equal(t, onion.PrivateKey(), onion2.PrivateKey())
		assert.Equal(t, onion.PubKey(), onion2.PubKey())
		assert.Equal(t, onion.Address(), onion2.Address())
	})

	t.Run("Save", func(t *testing.T) {
		// Generate new Onion
		onion, err := New()
		assert.NoError(t, err)

		// Save to session file
		err = onion.Save()
		assert.NoError(t, err)

		// Check if the session file was created
		filename := SESSION_DIR + "/" + onion.Session()
		_, err = os.Stat(filename)
		assert.NoError(t, err)

		// Check if the content of the session file is correct
		content, err := os.ReadFile(filename)
		assert.NoError(t, err)
		assert.Equal(t, onion.PrivateKey(), content)

		// Clean up session file
		err = os.Remove(filename)
		assert.NoError(t, err)

		// Cleanup folder
		err = os.Remove(SESSION_DIR)
		assert.NoError(t, err)
	})
}
