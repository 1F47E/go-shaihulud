package listner

import (
	"bufio"
	"context"
	"fmt"
	"go-dmtor/pkg/client/connection"
	"go-dmtor/pkg/client/message"
	cfg "go-dmtor/pkg/config"
	"go-dmtor/pkg/interfaces"
	"go-dmtor/pkg/logger"
	"io"
	"time"
)

var log = logger.New()

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
	log.Debug("Listner.Sender: Starting")
	defer func() {
		log.Warn("Sender: exit")
		// l.disconnect()
		//c.cancel() // do not cancel the main app ctx
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

// goroutine per connection
func (l *Listner) Receiver(user *connection.Connection, crypter interfaces.Asymmetric) {
	log.Debug("Listner.Receiver: Starting")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		log.Warn("Listner: exit")
		// c.disconnect()
		ticker.Stop()
		// c.cancel() // to not exit the app on disconnect, wait for another connection
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
					// The connection was closed.
					// msg := ""
					// if l.conn.Name != "" {
					// 	msg = fmt.Sprintf("<%s> disconnected", l.conn.Name)
					// } else {
					// 	msg = "Client disconnected"
					// }
					// log.Warn(msg)
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

			// Messages
			switch msg.Type {

			case message.HELLO:
				log.Info(">>Hello!")

			case message.ACK:
				log.Info(">>Ack!")

			case message.MSG:
				log.Debugf("\nraw msg %d, data %d bytes:\n=====\n%x\n=====\n", len(msg.Body), len(msg.Data()), msg.Body)
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
