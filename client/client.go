package client

import (
	"context"
	"fmt"
	"net"
	"time"

	cfg "go-dmtor/config"
	"go-dmtor/logger"
)

var log = logger.New()

type Client struct {
	conn  net.Conn
	ctx   context.Context
	addr  string
	MsgCh chan string
}

func NewClient(addr string) *Client {
	return &Client{
		ctx:   context.Background(),
		addr:  addr,
		MsgCh: make(chan string),
	}
}

func (c *Client) close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) ServerStart() {
	// open tcp connection to port 3000 and listen to incoming connections.
	// on connection print hello
	addr, err := net.ResolveTCPAddr("tcp4", c.addr)
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

func (c *Client) sender() {
	defer func() {
		log.Info("Sender: Closing connection")
		c.close()
	}()
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.MsgCh:
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
			bytes := make([]byte, cfg.MSG_MAX_SIZE)
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
