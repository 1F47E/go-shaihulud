package client_local

import (
	"context"
	"fmt"
	"net"

	"github.com/1F47E/go-shaihulud/internal/client/message"
)

type ClientLocal struct {
	ctx    context.Context
	cancel context.CancelFunc
	msgCh  chan message.Message
}

func New(ctx context.Context, cancel context.CancelFunc, msgCh chan message.Message) *ClientLocal {
	return &ClientLocal{
		ctx:    ctx,
		cancel: cancel,
		msgCh:  msgCh,
	}
}

func (c *ClientLocal) RunServer(port int, _ []byte) (net.Listener, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (c *ClientLocal) RunClient(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("RunClient connection error: %w", err)
	}
	return conn, err
}
