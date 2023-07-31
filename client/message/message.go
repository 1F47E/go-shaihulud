/*
MESSAGE FORMAT
-----------------
| 4 bytes
| message type
-----------------
| 4 bytes
| message nonce
-----------------
| 4 bytes
| message len
-----------------
| len bytes
| bytes body
-----------------
*/

package message

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"go-dmtor/logger"
	mrand "math/rand"
)

var log = logger.New()

type MsgType uint32

const (
	HELLO MsgType = iota
	ACK
	MSG
	KEY
	RUOK
	IMOK
	DISCONNECT
)

type Message struct {
	Type  MsgType
	Nonce uint32 // for ack
	Len   uint32
	Body  []byte
}

func NewMSG(msg []byte) Message {
	return Message{
		Type:  MSG,
		Nonce: nonce(),
		Len:   uint32(len(msg)),
		Body:  msg,
	}
}

func NewAck() Message {
	return Message{
		Type:  ACK,
		Nonce: nonce(),
		Len:   0,
		Body:  nil,
	}
}

func NewHello() Message {
	return Message{
		Type:  HELLO,
		Nonce: nonce(),
		Len:   0,
		Body:  nil,
	}
}

func NewKey(key []byte) Message {
	return Message{
		Type:  KEY,
		Nonce: nonce(),
		Len:   uint32(len(key)),
		Body:  key,
	}
}

func nonce() uint32 {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return mrand.Uint32()
	}
	return binary.BigEndian.Uint32(b)
}

func (t MsgType) String() string {
	switch t {
	case HELLO:
		return "HELLO"
	case ACK:
		return "ACK"
	case MSG:
		return "MSG"
	case KEY:
		return "KEY"
	case DISCONNECT:
		return "DISCONNECT"
	case RUOK:
		return "RUOK"
	case IMOK:
		return "IMOK"
	default:
		return "Unknown"
	}
}

func (m *Message) Data() []byte {
	if m.Len == 0 {
		return nil
	}
	if m.Len > uint32(len(m.Body)) {
		return m.Body
	}
	return m.Body[:m.Len]
}

func (m *Message) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// msg type
	if err := binary.Write(buf, binary.BigEndian, m.Type); err != nil {
		return nil, err
	}

	// msg nonce
	if err := binary.Write(buf, binary.BigEndian, m.Nonce); err != nil {
		return nil, err
	}

	// msg len
	if err := binary.Write(buf, binary.BigEndian, m.Len); err != nil {
		return nil, err
	}

	// msg body
	if err := binary.Write(buf, binary.BigEndian, m.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize(data []byte) (*Message, error) {
	buf := bytes.NewReader(data)
	// log message bytes
	log.Debugf("message bytes: %v\n", data)

	// type
	var msgType uint32
	if err := binary.Read(buf, binary.BigEndian, &msgType); err != nil {
		log.Errorf("error reading message type: %v\n", err)
		return nil, err
	}
	log.Debugf("Deserialized message type: %v\n", msgType)

	// msg nonce
	var msgNonce uint32
	if err := binary.Read(buf, binary.BigEndian, &msgNonce); err != nil {
		log.Errorf("error reading message nonce: %v\n", err)
		return nil, err
	}

	// msg len
	var msgLen uint32
	if err := binary.Read(buf, binary.BigEndian, &msgLen); err != nil {
		log.Errorf("error reading message len: %v\n", err)
		return nil, err
	}
	log.Debugf("Deserialized message len: %v\n", msgLen)

	// msg body
	body := make([]byte, msgLen)
	if err := binary.Read(buf, binary.BigEndian, &body); err != nil {
		log.Errorf("error reading message body: %v\n", err)
		return nil, err
	}
	log.Debugf("Deserialized message body: %v\n", body)

	return &Message{MsgType(msgType), msgNonce, msgLen, body}, nil
}
