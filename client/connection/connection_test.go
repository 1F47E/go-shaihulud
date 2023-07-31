package connection

import (
	"crypto/sha256"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnection_Update(t *testing.T) {
	t.Run("update public key and check name", func(t *testing.T) {
		// Mock connection
		mockConn := new(net.TCPConn) // You may need to replace this with an actual mock

		// Create a new connection
		conn := New(mockConn)

		// Check that handshake is not done
		assert.False(t, conn.Handshaked())

		// Mock public key
		pubKey := []byte("mockPublicKey")

		// Update the connection with the mock public key
		conn.Updade(pubKey)

		// Check that handshake is now done
		assert.True(t, conn.Handshaked())
		assert.Equal(t, pubKey, conn.PubKey)

		// Calculate expected name (first 2 bytes of SHA256 hash of public key, in upper case)
		sha256Hash := sha256.Sum256(pubKey)
		expectedName := fmt.Sprintf("%X", sha256Hash[:2])

		// Check that name is correctly set
		assert.Equal(t, expectedName, conn.Name)
	})
}
