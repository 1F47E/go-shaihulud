package client

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net"
	"time"

	"go-dmtor/client/message"
	cfg "go-dmtor/config"
	"go-dmtor/crypto"
	"go-dmtor/logger"
)

var log = logger.New()

type Client struct {
	conn        net.Conn
	ctx         context.Context
	cancel      context.CancelFunc
	addr        string
	msgCh       chan message.Message
	privKey     rsa.PrivateKey
	pubKey      rsa.PublicKey
	guestPubKey rsa.PublicKey
}

func NewClient(ctx context.Context, cancel context.CancelFunc, addr string) *Client {
	key := crypto.Keygen()
	return &Client{
		ctx:     ctx,
		cancel:  cancel,
		addr:    addr,
		msgCh:   make(chan message.Message),
		privKey: key,
		pubKey:  key.PublicKey,
	}
}

func (c *Client) close() {
	if c.conn != nil {
		c.conn.Close()
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
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Errorf("listen error: %v\n", err)
		return err
	}

	// TODO: cancel on ctx.Done()
	// BUG: blocking on accept

	log.Infof("Listening on %s", addr)
	c.conn, err = listener.Accept()
	if err != nil {
		log.Errorf("accept error: %v\n", err)
		return err
	}

	// TODO: send hello
	// w, err := c.conn.Write([]byte("hello"))
	// if err != nil {
	// 	log.Fatalf("write error: %v\n", err)
	// }
	// log.Printf("Wrote %d bytes\n", w)
	// TODO: do handshake

	go c.listner()
	go c.sender()
	return nil
}

func (c *Client) ServerConnect() error {
	var err error
	// connect to tcp port 3000, send user input
	c.conn, err = net.Dial("tcp", c.addr)
	if err != nil {
		log.Errorf("dial error: %v\n", err)
		return err
	}
	log.Infof("Connected to %s\n", c.addr)

	go c.sender()
	go c.listner()

	return nil
}

func (c *Client) SendMessage(msg []byte) {
	inputCipher := crypto.Encrypt(msg, &c.guestPubKey)
	log.Debugf("inputCipher: %d %v\n", len(inputCipher), inputCipher)
	c.msgCh <- message.NewMSG(inputCipher)
}

func (c *Client) sender() {
	defer func() {
		log.Info("Sender: Closing connection")
		c.close()
		c.cancel()
	}()

	// do handshake
	go func() {
		// send hello
		// c.msgCh <- message.NewHello()
		// send our PEM key
		pem, err := crypto.EncodePublicKeyToBytes(&c.pubKey)
		if err != nil {
			log.Errorf("Sender: PEM pub key error: %v\n", err)
			return
		}
		c.msgCh <- message.NewKey(pem)
		log.Debug("Sender: sent key")
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.msgCh:
			// send bytes to the connection
			mBytes, _ := msg.Serialize()
			w, err := c.conn.Write(mBytes)
			if err != nil {
				log.Fatalf("write error: %v\n", err)
			}
			log.Debugf("Wrote %d bytes\n", w)
		}
	}
}

func (c *Client) listner() {
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		log.Info("Listner: Closing connection")
		c.close()
		ticker.Stop()
		c.cancel()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			// read bytes from the connection
			bytes := make([]byte, cfg.MSG_MAX_SIZE)
			n, err := c.conn.Read(bytes)
			if err != nil {
				log.Errorf("Listner: read error: %v\n", err)
				// TODO: do not disconnect, wait for reconnect
				// continue
				return
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
				guestPubKey, err := crypto.DecodePublicKeyFromBytes(msg.Body)
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
