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
		// 270 bytes
		pubKey := []byte("AF3EAFDE09FBA80741034641180F13E029B056BC5F7440598EAC2EBFFE894D6C51D5263782D957FC95A856E1469159BFC97228448D2BF5F2DC896CE25758EF742235A7CEA5032C3F0B0B8A78EB8B08BA7D036E436F563078E660ED46")

		// Update the connection with the mock public key
		err := conn.UpdadeKey(pubKey)
		assert.NoError(t, err)

		// Check that handshake is now done
		assert.True(t, conn.Handshaked())
		assert.Equal(t, pubKey, conn.PubKey)

		// Calculate expected name (first 2 bytes of SHA256 hash of public key, in upper case)
		conn.UpdateName()
		sha256Hash := sha256.Sum256(conn.PubKey)
		expectedName := fmt.Sprintf("%X", sha256Hash[0:2])

		// Check that name is correctly set
		assert.Equal(t, expectedName, conn.Name)
	})
}
