package client_tor

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/1F47E/go-shaihulud/client/message"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
)

type TorClient struct {
	ctx    context.Context
	cancel context.CancelFunc
	msgCh  chan message.Message
}

func New(ctx context.Context, cancel context.CancelFunc, msgCh chan message.Message) *TorClient {
	return &TorClient{
		ctx:    ctx,
		cancel: cancel,
		msgCh:  msgCh,
	}
}

// Start the tor service and return the listener
func (c *TorClient) RunServer(_ string, onionPrivKey []byte) (net.Listener, error) {
	t, err := tor.Start(c.ctx, nil)
	if err != nil {
		return nil, err
	}
	// get key from bytes
	keyPair := ed25519.PrivateKey(onionPrivKey)
	torconn, err := t.Listen(c.ctx, &tor.ListenConf{Key: keyPair, LocalPort: 3000, RemotePorts: []int{80}})
	if err != nil {
		return nil, err
	}

	return torconn, nil
}

func (c *TorClient) RunClient(address string) (net.Conn, error) {
	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	t, err := tor.Start(c.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("tor start error: %w", err)
	}

	// Wait at most a minute to start network and get
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	// custom tor dialer
	dialer, err := t.Dialer(dialCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("tor dialer create error: %w", err)
	}

	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		if strings.Contains(err.Error(), "host unreachable") {
			return nil, fmt.Errorf("tor server is down :(")
		}
		return nil, fmt.Errorf("tor dial error: %w", err)
	}

	return conn, nil
}
