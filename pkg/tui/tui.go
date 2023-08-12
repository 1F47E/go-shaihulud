package tui

import (
	"context"
)

type TUI struct {
	ctx      context.Context
	eventsCh chan Event
}

func New(ctx context.Context, eventsCh chan Event) *TUI {
	return &TUI{ctx, eventsCh}
}

func (t *TUI) Run() {

	widget := NewWidget()
	// read events from channel and update spinner/progress bar
	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return

			case event := <-t.eventsCh:
				switch event.eventType {
				case eventTypeSpin:
					widget.SetSpinner(event.text)
				case eventTypeBar:
					widget.SetProgress(event.text, event.percent)
				case eventTypeText:
					widget.SetText(event.text)
				case eventTypeAccess:
					widget.SetAccess(event.access.key, event.access.password)
				}
			}
		}
	}()
	// init bubbletea spinner and progress bar
	widget.Run()
}
