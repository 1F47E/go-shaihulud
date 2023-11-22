package listner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/1F47E/go-shaihulud/internal/client/connection"
	"github.com/1F47E/go-shaihulud/internal/client/message"
	cfg "github.com/1F47E/go-shaihulud/internal/config"
	"github.com/1F47E/go-shaihulud/internal/interfaces"
	"github.com/1F47E/go-shaihulud/internal/logger"
)

type Listner struct {
	ctx    context.Context
	cancel context.CancelFunc
	msgCh  chan message.Message
}

func New(ctx context.Context, cancel context.CancelFunc, msgCh chan message.Message) *Listner {
	return &Listner{
		ctx:    ctx,
		cancel: cancel,
		msgCh:  msgCh,
	}
}

// goroutine per connection
func (l *Listner) Sender(user *connection.Connection, crypter interfaces.Asymmetric) {
	log := logger.New()
	log.Debug("Listner.Sender: Starting")
	defer func() {
		log.Debug("Sender: exit")
		l.cancel() // cancel only listners&senders ctx
	}()

	// do handshake
	// send our pubkic key
	go func() {
		l.msgCh <- message.NewKey(crypter.PubKey())
		log.Debug("Sender: sent key")
	}()

	// TODO: sign every message with a HMAC from password
	writer := bufio.NewWriter(user.Conn)
	for {
		select {
		case <-l.ctx.Done():
			return
		case msg := <-l.msgCh:
			// TODO: check is there was a handshake
			log.Debugf("Sender: Got msg: %v\n", msg)
			// send bytes to the connection
			mBytes, _ := msg.Serialize()

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

// goroutine per connection
func (l *Listner) Receiver(user *connection.Connection, crypter interfaces.Asymmetric) {
	log := logger.New()

	log.Debug("Listner.Receiver: Starting")

	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		log.Debug("Listner: exit")
		ticker.Stop()
		l.cancel() // cancel only local context for listners
	}()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			// read bytes from the connection
			if user.Conn == nil {
				log.Warn("Listner: No connection")
				return
			}
			bytes := make([]byte, cfg.MSG_MAX_SIZE)
			n, err := user.Conn.Read(bytes)
			if err != nil {
				if err == io.EOF {
					log.Warn("Listner: Connection closed")
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

			// react on incoming message

			switch msg.Type {

			case message.HLLO:
				log.Info(">>Hello!")

			case message.ACK:
				log.Debugf(">> Ack! msg %d delivered", msg.Nonce)
				println("☑︎")

			case message.MSG:
				log.Debugf("\nraw msg %d bytes:\n=====\n%x\n=====\n", len(msg.Body), msg.Body)
				// decode msg
				decrypted, err := crypter.Decrypt(msg.Body)
				if err != nil {
					log.Errorf("error decrypting msg: %v\n", err)
					continue
				}
				now := time.Now().Format("15:04:05")
				fmt.Printf("%s <%s> %s\n", now, user.Name, string(decrypted))

			case message.KEY:
				log.Debugf("got public key from user: %d bytes\n%v", len(msg.Body), msg.Body)
				// save guest key
				err := user.UpdadeKey(msg.Body)
				if err != nil {
					log.Error("got wrong pub key from the user")
					// TODO: disconnect user
				} else {
					user.UpdateName()
					log.Infof("<%s> entered the chat", user.Name)
				}
			case message.DISC:
				log.Warnf("<%s> disconnected", user.Name)

			default:
				log.Warnf("unknown message type: %s\n", msg.Type)
				if msg.Len > 0 {
					fmt.Printf("len: %d\n", msg.Len)
					fmt.Printf("data: %s\n", string(msg.Body))
				}
			}

			// send delivery confirmation (ACK)
			if msg.Type != message.ACK && msg.Nonce > 0 {
				l.msgCh <- message.NewAck(msg.Nonce)
			}
		}
	}
}
