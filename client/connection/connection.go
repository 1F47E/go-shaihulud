package connection

import (
	"crypto/rsa"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
)

type Connection struct {
	UUID string
	Conn net.Conn
	// ctx    context.Context
	// cancel context.CancelFunc
	Name   string
	PubKey rsa.PublicKey
}

func New(conn net.Conn) *Connection {
	// c, cancel := context.WithCancel(ctx)
	return &Connection{
		UUID: uuid.New().String(),
		Conn: conn,
		// ctx:    c,
		// cancel: cancel,
	}
}

func (c *Connection) Updade(key *rsa.PublicKey) {
	c.PubKey = *key
	// make a name from first bytes or the pub key
	hash := fmt.Sprintf("%x", key.N.Bytes()[0:2])
	hash = strings.ToUpper(hash)
	c.Name = hash
}
