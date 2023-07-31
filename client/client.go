package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"go-dmtor/client/connection"
	"go-dmtor/client/message"
	cfg "go-dmtor/config"
	is "go-dmtor/interfaces"
	"go-dmtor/logger"
)

var log = logger.New()

type Client struct {
	ctx     context.Context
	cancel  context.CancelFunc
	addr    string
	conn    *connection.Connection
	msgCh   chan message.Message
	crypter is.Asymmetric
}

func NewClient(ctx context.Context, cancel context.CancelFunc, addr string, crypter is.Asymmetric) *Client {
	return &Client{
		ctx:     ctx,
		cancel:  cancel,
		addr:    addr,
		msgCh:   make(chan message.Message),
		crypter: crypter,
	}
}

func (c *Client) disconnect() {
	if c.conn != nil && c.conn.Conn != nil {
		c.conn.Conn.Close()
	}
}

func (c *Client) ServerStart() error {
	defer func() {
		log.Warnf("ServerStart exit\n")
	}()
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
	log.Infof("Main listening socket is open at %s", addr)

	// run listener accept in a sep G to allow shutdown
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()

		default:
			log.Info("Waiting for connection...")
			// block
			conn, err := listener.Accept()
			if err != nil {
				log.Errorf("accept error: %v\n", err)
			}
			ip := conn.RemoteAddr().String()
			log.Debugf("Connection open for %s\n", ip)

			if c.conn != nil && c.conn.Name != "" {
				log.Warnf("<%s> reconnected\n", c.conn.Name)
				c.conn.Conn = conn
			} else {
				c.conn = connection.New(conn)
				log.Warn("Connecting...\n")
				log.Debugf("Connection uuid is [%s]\n", c.conn.UUID)
			}

			// custom ctx to cancel both listner and sender
			ctx, cancel := context.WithCancel(c.ctx)
			go c.listner(ctx, cancel)
			go c.sender(ctx, cancel)
			log.Debugf("ServerStart: blocing, working with connection")
			<-ctx.Done()
			// Block here untill we have a connection
		}
	}
}

func (c *Client) ServerConnect() error {
	defer func() {
		log.Warnf("ServerConnect exit\n")
	}()
	retry := 0
	isFirstTry := true
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
			conn, err := net.Dial("tcp", c.addr)
			if err != nil {
				if isFirstTry {
					return fmt.Errorf("Server no avaiable, retry later")
				}
				retry++
				if retry > cfg.CLIENT_MAX_RETRY {
					return fmt.Errorf("Max retry reached, Connection error: %v\n", err)
				}
				log.Warnf("Connection error: %v, retry %d/%d\n", err, retry, cfg.CLIENT_MAX_RETRY)
				// TODO: rewrite to ticker to cancel faster?
				time.Sleep(time.Duration(retry) * 1 * time.Second)
				continue
			}
			log.Infof("Connected to %s\n", c.addr)
			isFirstTry = false

			c.conn = connection.New(conn)
			ctx, cancel := context.WithCancel(c.ctx)
			go c.sender(ctx, cancel)
			go c.listner(ctx, cancel)
			<-ctx.Done()
		}
	}
}

func (c *Client) SendMessage(msg []byte) error {
	inputCipher, err := c.crypter.Encrypt(msg, c.conn.PubKey)
	if err != nil {
		return err
	}
	log.Debugf("inputCipher: %d %v\n", len(inputCipher), inputCipher)
	c.msgCh <- message.NewMSG(inputCipher)
	return nil
}

func (c *Client) sender(ctx context.Context, cancel context.CancelFunc) {
	defer func() {
		log.Info("Sender: Closing connection")
		c.disconnect()
		//c.cancel() // do not cancel the main app ctx
		cancel() // cancel only listners&senders ctx
	}()

	// do handshake
	go func() {
		// send our pubkic key
		pubKey := c.crypter.PubKey()
		c.msgCh <- message.NewKey(pubKey)
		log.Debug("Sender: sent key")
	}()

	writer := bufio.NewWriter(c.conn.Conn)
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.msgCh:
			// TODO: check is there was a handshake
			log.Debugf("Sender: Got msg: %v\n", msg)
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

func (c *Client) listner(ctx context.Context, cancel context.CancelFunc) {
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
			if c.conn == nil || c.conn.Conn == nil {
				log.Warn("Listner: No connection")
				return
			}
			bytes := make([]byte, cfg.MSG_MAX_SIZE)
			n, err := c.conn.Conn.Read(bytes)
			if err != nil {
				if err == io.EOF {
					// The connection was closed.
					msg := ""
					if c.conn.Name != "" {
						msg = fmt.Sprintf("<%s> disconnected", c.conn.Name)
					} else {
						msg = "Client disconnected"
					}
					log.Warn(msg)
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
				log.Debugf("raw msg:\n=====\n%s\n=====\n", string(msg.Body))
				// decode msg
				decrypted, err := c.crypter.Decrypt(msg.Data())
				if err != nil {
					log.Errorf("error decrypting msg: %v\n", err)
					continue
				}
				now := time.Now().Format("15:04:05")
				fmt.Printf("%s <%s> %s\n", now, c.conn.Name, string(decrypted))
			case message.KEY:
				log.Debugf("got public key from user: %d bytes\n%v", len(msg.Body), msg.Body)
				// save guest key
				c.conn.Updade(msg.Body)
				log.Infof("<%s> entered the chat", c.conn.Name)
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
