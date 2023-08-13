package tui

import (
	"context"

	tauth "github.com/1F47E/go-shaihulud/pkg/tui/t.auth"
	tstatus "github.com/1F47E/go-shaihulud/pkg/tui/t.status"
)

type TUI struct {
	ctx      context.Context
	eventsCh chan Event
	loader   *tstatus.StatusWidget
	auth     *tauth.AuthWidget
}

func New(ctx context.Context, eventsCh chan Event) *TUI {
	return &TUI{
		ctx:      ctx,
		eventsCh: eventsCh,
		loader:   tstatus.New(),
		auth:     tauth.New(),
	}
}

// renders are blocking
func (t *TUI) RenderAuth() (string, string, error) {
	key, password, err := t.auth.Run()
	if err != nil {
		return "", "", err
	}
	return key, password, nil
}

func (t *TUI) RenderLoader() {
	t.loader.SetText("")
	t.loader.Run()
}

// read events from channel and update spinner/progress bar
func (t *TUI) Listner() {
	for {
		select {
		case <-t.ctx.Done():
			return

		case event := <-t.eventsCh:
			switch event.eventType {
			case eventTypeSpin:
				t.loader.SetSpinner(event.text)
			case eventTypeBar:
				t.loader.SetProgress(event.text, event.percent)
			case eventTypeText:
				t.loader.SetText(event.text)
			case eventTypeAccess:
				t.loader.SetAccess(event.access.key, event.access.password)
			}
		}
	}
}
