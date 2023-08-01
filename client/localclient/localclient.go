package localclient

import (
	"context"
	"fmt"
	"go-dmtor/client/message"
	"net"
)

type LocalClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	msgCh  chan message.Message
}

func New(ctx context.Context, cancel context.CancelFunc, msgCh chan message.Message) *LocalClient {
	return &LocalClient{
		ctx:    ctx,
		cancel: cancel,
		msgCh:  msgCh,
	}
}

func (c *LocalClient) RunServer(address string, _ []byte) (net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (c *LocalClient) RunClient(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("RunClient connection error: %w", err)
	}
	return conn, err
}
