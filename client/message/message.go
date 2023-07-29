/*
MESSAGE FORMAT
-----------------
| 4 bytes
| message type
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
	"encoding/binary"
	"go-dmtor/logger"
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
	Type MsgType
	Len  uint32
	Body []byte
}

func NewMessageText(msg string) *Message {
	return &Message{
		Type: MSG,
		Len:  uint32(len(msg)),
		Body: []byte(msg),
	}
}

func NewMessageHello() *Message {
	return &Message{
		Type: HELLO,
		Len:  0,
		Body: nil,
	}
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

func (m *Message) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, m.Type); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, m.Len); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, m.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize(data []byte) (*Message, error) {
	buf := bytes.NewReader(data)
	// log message bytes
	log.Debugf("message bytes: %v\n", data)

	var msgType uint32
	if err := binary.Read(buf, binary.BigEndian, &msgType); err != nil {
		log.Errorf("error reading message type: %v\n", err)
		return nil, err
	}
	log.Debugf("Deserialized message type: %v\n", msgType)

	var msgLen uint32
	if err := binary.Read(buf, binary.BigEndian, &msgLen); err != nil {
		log.Errorf("error reading message len: %v\n", err)
		return nil, err
	}
	log.Debugf("Deserialized message len: %v\n", msgLen)

	body := make([]byte, 0)
	if msgLen > 0 {
		body = make([]byte, msgLen)
		if err := binary.Read(buf, binary.BigEndian, &body); err != nil {
			log.Errorf("error reading message body: %v\n", err)
			return nil, err
		}
	}
	log.Debugf("Deserialized message body: %v\n", body)

	return &Message{MsgType(msgType), msgLen, body}, nil
}
