package client

import (
	"bufio"
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"time"

	"go-dmtor/client/connection"
	"go-dmtor/client/message"
	cfg "go-dmtor/config"
	"go-dmtor/crypto"
	"go-dmtor/logger"
)

var log = logger.New()

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	addr   string
	// conn        net.Conn
	connections map[uint64]*connection.Connection
	msgCh       chan message.Message

	// TODO: make a keychain for all the users in a chat
	privKey     rsa.PrivateKey
	pubKey      rsa.PublicKey
	guestPubKey rsa.PublicKey
}

func NewClient(ctx context.Context, cancel context.CancelFunc, addr string) *Client {
	key := crypto.Keygen()
	return &Client{
		ctx:         ctx,
		cancel:      cancel,
		addr:        addr,
		connections: make(map[uint64]*connection.Connection),
		msgCh:       make(chan message.Message),
		privKey:     key,
		pubKey:      key.PublicKey,
	}
}

func (c *Client) disconnect() {
	// close all connections
	for _, conn := range c.connections {
		if conn != nil && conn.Conn != nil {
			conn.Conn.Close()
		}
	}
}

func (c *Client) ServerStart() error {
	// open tcp connection to port 3000 and listen to incoming connections.
	// on connection print hello
	addr, err := net.ResolveTCPAddr("tcp4", c.addr)
	if err != nil {
		log.Errorf("resolve error: %v\n", err)
		return err
	}

	log.Infof("Listening on %s", addr)
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("listen error: %v\n", err)
		return err
	}

	// TODO: handle multiple connections
	// run listener accept in a sep G to allow shutdown
	go func() {
		defer func() {
			log.Warnf("ServerStart exit\n")
		}()
		for {
			select {
			case <-c.ctx.Done():
				return

			default:
				log.Infof("Listening on %s", addr)
				conn, err := listener.Accept()
				if err != nil {
					log.Errorf("accept error: %v\n", err)
				}
				// connID := uuid.New().String()
				ip := conn.RemoteAddr().String()
				connID := crypto.Hash([]byte(ip))

				if _, ok := c.connections[connID]; ok {
					log.Info("Client reconnected: %s\n", connID)
					c.connections[connID].Conn = conn
				} else {
					c.connections[connID] = &connection.Connection{
						ID:   connID,
						Conn: conn,
					}
					log.Infof("Accepted connection from %s\n", connID)
				}
				log.Debugf("Total connections: %d\n", len(c.connections))

				// TODO: WIP pass connection to listners

				// custom ctx to cancel both listner and sender
				ctx, cancel := context.WithCancel(c.ctx)
				go c.listner(conn, ctx, cancel)
				go c.sender(conn, ctx, cancel)
				log.Warn("ServerStart: blocing, working with connection")
				<-ctx.Done()
				// Block here untill we have a connection and listners are running
			}
		}
	}()

	return nil
}

func (c *Client) ServerConnect() error {
	var err error
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Errorf("dial error: %v\n", err)
		return err
	}
	log.Infof("Connected to %s\n", c.addr)

	ctx, cancel := context.WithCancel(c.ctx)
	go c.sender(conn, ctx, cancel)
	go c.listner(conn, ctx, cancel)
	// TODO: make reconnect here

	return nil
}

func (c *Client) SendMessage(msg []byte) {
	inputCipher := crypto.Encrypt(msg, &c.guestPubKey)
	log.Debugf("inputCipher: %d %v\n", len(inputCipher), inputCipher)
	c.msgCh <- message.NewMSG(inputCipher)
}

func (c *Client) sender(conn net.Conn, ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		log.Info("Sender: Closing connection")
		c.disconnect()
		//c.cancel() // do not cancel the main app ctx
		cancel() // cancel only listners&senders ctx
	}()

	// do handshake
	go func() {
		// send hello
		// c.msgCh <- message.NewHello()
		// send our PEM key
		pem, err := crypto.PubToBytes(&c.pubKey)
		if err != nil {
			log.Errorf("Sender: PEM pub key error: %v\n", err)
			return
		}
		c.msgCh <- message.NewKey(pem)
		log.Debug("Sender: sent key")
	}()
	writer := bufio.NewWriter(conn)
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.msgCh:
			// send bytes to the connection
			mBytes, _ := msg.Serialize()
			// w, err := conn.Write(mBytes)
			// if err != nil {
			// 	log.Fatalf("write error: %v\n", err)
			// }
			w, err := writer.Write(mBytes)
			if err != nil {
				log.Fatalf("Sender: write error: %v\n", err)
			}
			err = writer.Flush()
			if err != nil {
				log.Errorf("Sender: Flush error: %v", err)
				return
			}
			log.Debugf("Sender: Wrote %d bytes\n", w)
		}
	}
}

func (c *Client) listner(conn net.Conn, ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		log.Info("Listner: Closing connection")
		c.disconnect()
		ticker.Stop()
		// c.cancel() // to not exit the app on disconnect, wait for another connection
		cancel() // cancel only local context for listners
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// read bytes from the connection
			bytes := make([]byte, cfg.MSG_MAX_SIZE)
			n, err := conn.Read(bytes)
			if err != nil {
				if err == io.EOF {
					// The connection was closed.
					log.Warnf("Listner: Connection closed: %v", err)
					return
				}
				log.Errorf("Listner: Read error: %v", err)
				continue
			}
			log.Debugf("Received: %d bytes:\n", n)
			log.Debugf("raw: %v\n", bytes)
			msg, err := message.Deserialize(bytes)
			if err != nil {
				log.Errorf("deserialize error: %v\n", err)
			}
			log.Debugf("Msg type: %s\n", msg.Type)
			switch msg.Type {
			case message.HELLO:
				log.Info(">>Hello!")
			case message.ACK:
				log.Info(">>Ack!")
			case message.MSG:
				log.Infof("raw msg: %s\n", string(msg.Body))
				// decode msg
				decrypted := crypto.Decrypt(msg.Data(), &c.privKey)
				log.Infof("decrypted: %s\n", string(decrypted))
			case message.KEY:
				log.Info(">>KEY")
				// decode guest key from bytes
				guestPubKey, err := crypto.BytesToPub(msg.Body)
				if err != nil {
					log.Errorf("decode key error: %v\n", err)
				} else {
					c.guestPubKey = *guestPubKey
					log.Infof("guest key received: %d - ", len(c.guestPubKey.N.Bytes()))
					log.Infof("(%x)", c.guestPubKey.N.Bytes()[0:3])
				}
			default:
				fmt.Printf(">>%s\n", msg.Type)
				if msg.Len > 0 {
					fmt.Printf("len: %d\n", msg.Len)
					fmt.Printf("data: %s\n", string(msg.Body))
				}
			}
			// TODO: send ack
		}
	}
}
