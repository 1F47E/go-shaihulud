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

	zlog "github.com/rs/zerolog/log"
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
	zlog.Debug().Msg("Listner.Sender: Starting")
	defer func() {
		zlog.Debug().Msg("Sender: exit")
		l.cancel() // cancel only listners&senders ctx
	}()

	// do handshake
	// send our pubkic key
	go func() {
		l.msgCh <- message.NewKey(crypter.PubKey())
		zlog.Debug().Msg("Sender: sent key")
	}()

	// TODO: sign every message with a HMAC from password
	writer := bufio.NewWriter(user.Conn)
	for {
		select {
		case <-l.ctx.Done():
			return
		case msg := <-l.msgCh:
			// TODO: check is there was a handshake
			zlog.Debug().Msgf("Sender: Got msg: %v\n", msg)
			// send bytes to the connection
			mBytes, _ := msg.Serialize()

			w, err := writer.Write(mBytes)
			if err != nil {
				zlog.Fatal().Err(err).Msg("Sender: write error")
			}
			err = writer.Flush()
			if err != nil {
				zlog.Error().Err(err).Msg("Sender: Flush error")
				return
			}
			zlog.Debug().Msgf("Sender: Wrote %d bytes\n", w)
		}
	}
}

// goroutine per connection
func (l *Listner) Receiver(user *connection.Connection, crypter interfaces.Asymmetric) {

	zlog.Debug().Msg("Listner.Receiver: Starting")

	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		zlog.Debug().Msg("Listner: exit")
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
				zlog.Warn().Msg("Listner: No connection")
				return
			}
			bytes := make([]byte, cfg.MSG_MAX_SIZE)
			n, err := user.Conn.Read(bytes)
			if err != nil {
				if err == io.EOF {
					zlog.Warn().Msg("Listner: Connection closed")
					return
				}
				zlog.Error().Err(err).Msg("Listner: Read error")
				continue
			}
			zlog.Debug().Msgf("Received: %d bytes:\n", n)
			zlog.Debug().Msgf("raw: %v\n", bytes)
			msg, err := message.Deserialize(bytes)
			if err != nil {
				zlog.Error().Err(err).Msg("deserialize error")
			}
			zlog.Debug().Msgf("Msg type: %s\n", msg.Type)

			// react on incoming message

			switch msg.Type {

			case message.HLLO:
				zlog.Info().Msg(">>Hello!")

			case message.ACK:
				zlog.Debug().Msgf(">> Ack! msg %d delivered", msg.Nonce)
				println("☑︎")

			case message.MSG:
				zlog.Debug().Msgf("\nraw msg %d bytes:\n=====\n%x\n=====\n", len(msg.Body), msg.Body)
				// decode msg
				decrypted, err := crypter.Decrypt(msg.Body)
				if err != nil {
					zlog.Error().Err(err).Msg("error decrypting msg")
					continue
				}
				now := time.Now().Format("15:04:05")
				fmt.Printf("%s <%s> %s\n", now, user.Name, string(decrypted))

			case message.KEY:
				zlog.Debug().Msgf("got public key from user: %d bytes\n%v", len(msg.Body), msg.Body)
				// save guest key
				err := user.UpdadeKey(msg.Body)
				if err != nil {
					zlog.Error().Msg("got wrong pub key from the user")
					// TODO: disconnect user
				} else {
					user.UpdateName()
					zlog.Info().Msgf("<%s> entered the chat", user.Name)
				}
			case message.DISC:
				zlog.Warn().Msgf("<%s> disconnected", user.Name)

			default:
				zlog.Warn().Msgf("unknown message type: %s\n", msg.Type)
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
