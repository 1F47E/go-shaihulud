package connection

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strings"

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

func (c *Connection) Updade(pubKey []byte) {
	c.PubKey = pubKey
	// make a name from first bytes or the hash of pub key
	sha256Hash := sha256.Sum256(pubKey)
	hash := fmt.Sprintf("%x", sha256Hash[0:2])
	hash = strings.ToUpper(hash)
	c.Name = hash
}
