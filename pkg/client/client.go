package client

import (
	"context"
	"net"
	"os"
	"strings"

	"github.com/1F47E/go-shaihulud/pkg/client/connection"
	"github.com/1F47E/go-shaihulud/pkg/client/listner"
	"github.com/1F47E/go-shaihulud/pkg/client/localclient"
	"github.com/1F47E/go-shaihulud/pkg/client/message"
	"github.com/1F47E/go-shaihulud/pkg/client/torclient"
	cfg "github.com/1F47E/go-shaihulud/pkg/config"
	myaes "github.com/1F47E/go-shaihulud/pkg/cryptotools/aes"
	"github.com/1F47E/go-shaihulud/pkg/cryptotools/auth"
	"github.com/1F47E/go-shaihulud/pkg/interfaces"
	"github.com/1F47E/go-shaihulud/pkg/logger"
)

var log = logger.New()

// can be local or tor
type Connector interface {
	RunServer(address string, onionPrivKey []byte) (net.Listener, error)
	RunClient(address string) (net.Conn, error)
}

type ConnectionType int

const (
	Local ConnectionType = iota
	Tor
)

type Client struct {
	ctx       context.Context
	cancel    context.CancelFunc
	msgCh     chan message.Message
	connector Connector
	crypter   interfaces.Asymmetric
	user      *connection.Connection
	listner   *listner.Listner
	connType  ConnectionType
}

func NewClient(ctx context.Context, cancel context.CancelFunc, connType ConnectionType, crypter interfaces.Asymmetric) *Client {
	msgCh := make(chan message.Message)
	var connector Connector

	// init connector debug or tor
	switch connType {
	case Local:
		connector = localclient.New(ctx, cancel, msgCh)
	case Tor:
		connector = torclient.New(ctx, cancel, msgCh)
	}

	// create listner
	lCtx, lCancel := context.WithCancel(ctx)
	lstnr := listner.New(lCtx, lCancel, msgCh)

	return &Client{
		ctx:       ctx,
		cancel:    cancel,
		msgCh:     msgCh,
		connector: connector,
		crypter:   crypter,
		listner:   lstnr,
		connType:  connType,
	}
}

func (c *Client) RunServer(session string) error {

	// generate auth key and password
	crypter := myaes.New()
	auth, err := auth.New(crypter, session)
	if err != nil {
		log.Fatalf("cant create auth: %v\n", err)
	}
	log.Warn("========== AUTH ==========")
	log.Warnf("Key: %s\n\n", auth.AccessKey())
	log.Warnf("Password: %s\n", auth.Password())
	log.Warn("==========================")

	// get address
	address := ""
	switch c.connType {
	case Local:
		log.Info("Starting local server...")
		address = "localhost:3000"
	case Tor:
		log.Info("Starting tor...")
		address = auth.OnionAddressFull()
		log.Debugf("onion address: %v\n", address)
	default:
		log.Fatalf("unknown connection type: %v\n", c.connType)
	}

	// run server with a given address
	log.Debugf("Client.RunServer: %v\n", address)
	listener, err := c.connector.RunServer(address, auth.Onion().PrivKey())
	if err != nil {
		return err
	}
	log.Info("Server started")

	// accept incoming connections

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				log.Debugf("RunServer listner context done: %v\n", c.ctx.Err())
				return

			default:
				log.Debug("Client.RunServer: Waiting for a connection")

				conn, err := listener.Accept()
				if err != nil {
					log.Errorf("Client.RunServer listener.Accept error: %v\n", err)
					continue
				}
				user := connection.New(conn) // connection with user data
				c.user = user
				log.Debug("Client.RunServer: Got a connection")

				// Create a new Listner for each connection
				ctx, cancel := context.WithCancel(c.ctx)
				listner := listner.New(ctx, cancel, c.msgCh)
				go listner.Sender(user, c.crypter)
				go listner.Receiver(user, c.crypter)
				go c.ListenUserInput()
			}
		}
	}()

	return nil
}

func (c *Client) RunClient(key, password string) error {

	// create auth struct and try to decode key
	aes := myaes.New()
	ath, err := auth.NewFromKey(aes, key, password)
	if err != nil {
		if strings.Contains(err.Error(), "authentication failed") {
			log.Fatal("wrong password")
		}
		log.Fatalf("cant create auth: %v\n", err)
	}

	// ===== At this point access key and pass are valid

	// get address to connect to
	address := ""
	switch c.connType {
	case Local:
		address = "localhost:3000"
		log.Infof("Connecting to %s...", address)
	case Tor:
		address = ath.OnionAddressFull()
		log.Info("Connecting to tor...")
		log.Debugf("onion address: %v\n", address)
	default:
		log.Fatalf("unknown connection type: %v\n", c.connType)
	}

	// Run the connector
	conn, err := c.connector.RunClient(address)
	if err != nil {
		return err
	}
	user := connection.New(conn) // connection with user data
	c.user = user

	// Run the listener, sender, and input listener goroutines
	ctx, cancel := context.WithCancel(c.ctx)
	c.listner = listner.New(ctx, cancel, c.msgCh)
	go c.listner.Sender(user, c.crypter)
	go c.listner.Receiver(user, c.crypter)
	go c.ListenUserInput()

	return nil
}

func (c *Client) ListenUserInput() {
	for {
		select {
		case <-c.ctx.Done():
			log.Warnf("context done: %v\n", c.ctx.Err())
			return
		default:
			input := make([]byte, cfg.MSG_MAX_SIZE)
			n, err := os.Stdin.Read(input)
			if err != nil {
				log.Fatalf("read error: %v\n", err)
				return
			}
			text := input[:n]
			log.Debugf("user input: %d %v\n", len(text), text)
			log.Debugf("crypter: %v\n", c.crypter)
			inputCipher, err := c.crypter.Encrypt(text, c.user.PubKey)
			if err != nil {
				log.Errorf("can't send a message: %v\n", err)
			}
			log.Debugf("inputCipher: %d %v\n", len(inputCipher), inputCipher)
			c.msgCh <- message.NewMSG(inputCipher)
		}
	}
}

func (c *Client) Close() {
	if c.user != nil && c.user.Conn != nil {
		c.user.Conn.Close()
	}
}
