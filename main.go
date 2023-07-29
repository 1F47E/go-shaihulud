package main

import (
	"context"
	"fmt"
	"go-dmtor/logger"
	"net"
	"os"
	"time"
)

var log = logger.New()

var msgCh = make(chan string)
var msgMaxSize = 1024

var addr = "localhost:3000"

var usage = "Usage: %s <srv|cli>\n"

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatalf(usage, args[0])
	}
	arg := args[1]
	if arg != "srv" && arg != "cli" {
		log.Fatalf(usage, args[0])
	}
	cli := NewClient()
	if arg == "srv" {
		cli.serverStart()
		// crypt_demo()
	} else {
		// read user input
		err := cli.serverConnect()
		if err != nil {
			log.Fatalf("connect error: %v\n", err)
		}
	}

	// block and wait for user input
	for {
		input := make([]byte, msgMaxSize)
		_, err := os.Stdin.Read(input)
		if err != nil {
			log.Fatalf("read error: %v\n", err)
			return
		}
		msgCh <- string(input)
	}
}

// ====== CLIENT

type Client struct {
	conn net.Conn
	ctx  context.Context
}

func NewClient() *Client {
	return &Client{
		ctx: context.Background(),
	}
}

func (c *Client) close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) serverStart() {
	// open tcp connection to port 3000 and listen to incoming connections.
	// on connection print hello
	addr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		log.Fatalf("resolve error: %v\n", err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	log.Infof("Listening on %s", addr)
	c.conn, err = listener.Accept()
	if err != nil {
		log.Fatalf("accept error: %v\n", err)
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
}

func (c *Client) serverConnect() error {
	var err error
	// connect to tcp port 3000, send user input
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		log.Errorf("dial error: %v\n", err)
		return err
	}
	log.Infof("Connected to %s\n", addr)

	go c.sender()
	go c.listner()

	return nil
}

func (c *Client) sender() {
	defer func() {
		log.Info("Sender: Closing connection")
		c.close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-msgCh:
			// send bytes to the connection
			w, err := c.conn.Write([]byte(msg))
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
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			// read bytes from the connection
			bytes := make([]byte, msgMaxSize)
			n, err := c.conn.Read(bytes)
			if err != nil {
				log.Errorf("read error: %v\n", err)
				return
			}
			log.Debugf("Received: %d bytes:\n", n)
			fmt.Printf("%s", bytes)
			// TODO: send ack
		}
	}
}

// func crypt_demo() {

// 	key := crypto.Keygen()
// 	// Get the public key
// 	publicKey := &key.PublicKey

// 	// Encrypt a message
// 	message := "hello, world"
// 	cipher := crypto.Encrypt(message, publicKey)
// 	fmt.Printf("Ciphertext: %x\n", cipher)

// 	plain := crypto.Decrypt(cipher, &key)
// 	fmt.Printf("Plaintext: %s\n", plain)
// }
