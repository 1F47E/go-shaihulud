package tui

type eventType int

const (
	eventTypeSpin eventType = iota
	eventTypeBar
	eventTypeText
	eventTypeTextError
	eventTypeAccess
	eventTypeChatReady
)

type access struct {
	key      string
	password string
}

type Event struct {
	eventType eventType
	text      string
	percent   float64
	access    access
}

func NewEventSpin(text string) Event {
	return Event{
		eventType: eventTypeSpin,
		text:      text,
	}
}

func NewEventBar(text string, percent float64) Event {
	return Event{
		eventType: eventTypeBar,
		text:      text,
		percent:   percent,
	}
}

func NewEventText(text string) Event {
	return Event{
		eventType: eventTypeText,
		text:      text,
	}
}

func NewEventError(text string) Event {
	return Event{
		eventType: eventTypeTextError,
		text:      text,
	}
}

func NewEventAccess(key, password string) Event {
	return Event{
		eventType: eventTypeAccess,
		access:    access{key, password},
	}
}
