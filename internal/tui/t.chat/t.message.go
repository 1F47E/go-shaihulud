package tchat

import (
	"fmt"
	"time"
)

type icon int

const (
	iconNone icon = iota
	iconSuccess
	iconError
	iconSent
)

type ChatMessage struct {
	isSelf   bool
	name     string
	text     string
	datetime time.Time
	icon     icon
}

func NewChatMessage(name, text string) *ChatMessage {
	return &ChatMessage{
		name:     name,
		isSelf:   name == "",
		text:     text,
		datetime: time.Now(),
		icon:     iconNone,
	}
}

func (m *ChatMessage) SetSent() {
	m.icon = iconSent
}

func (m *ChatMessage) SetSuccess() {
	m.icon = iconSuccess
}

func (m *ChatMessage) SetError() {
	m.icon = iconError
}

/*
messages = append(messages, senderStyle.Render("You:  ")+fmt.Sprintf("message %d", i))
messages = append(messages, senderStyle.Render("You: ")+styleRed.Render(" ⨯ ")+fmt.Sprintf("message %d", i))
messages = append(messages, senderStyle.Render("You: ")+styleYellow.Render(" … ")+fmt.Sprintf("message %d", i))
messages = append(messages, senderStyle.Render("You: ")+styleGreen.Render(" ☑︎ ")+fmt.Sprintf("message %d", i))
*/
func (d *ChatMessage) Text() string {
	icon := ""
	switch d.icon {
	case iconSuccess:
		icon = styleGreen.Render(" ☑︎ ")
	case iconError:
		icon = styleRed.Render(" ⨯ ")
	case iconSent:
		icon = styleYellow.Render(" … ")
	}
	t := d.datetime.Format("15:04:05")
	tStr := styleGray.Render(t)
	name := ""
	if d.isSelf {
		name = styleSender.Render("You")
	} else {
		name = styleReceiver.Render(d.name)
	}

	return fmt.Sprintf("%s %s\t%s\n\t\t%s", icon, tStr, name, d.text)
}
