package connection

import (
	"crypto/sha256"
	"fmt"
	"net"

	"github.com/google/uuid"
)

type Connection struct {
	UUID   string
	Conn   net.Conn
	Name   string
	PubKey []byte // if nil - no handshake yet
}

func New(conn net.Conn) *Connection {
	return &Connection{
		UUID: uuid.New().String(),
		Conn: conn,
	}
}

func (c *Connection) Handshaked() bool {
	return c.PubKey != nil
}

func (c *Connection) UpdadeKey(pubKey []byte) error {
	c.PubKey = pubKey
	return nil
}

func (c *Connection) UpdateName() {
	// make a name from first bytes or the hash of pub key
	sha256Hash := sha256.Sum256(c.PubKey)
	hash := fmt.Sprintf("%X", sha256Hash[0:2])
	c.Name = hash
}
