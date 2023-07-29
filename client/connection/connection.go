package connection

import "net"

type Connection struct {
	ID   uint64
	Conn net.Conn
}
