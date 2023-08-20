package tui

import (
	"context"
	"fmt"

	tauth "github.com/1F47E/go-shaihulud/tui/t.auth"
	tpage "github.com/1F47E/go-shaihulud/tui/t.chat"
	tstatus "github.com/1F47E/go-shaihulud/tui/t.status"
	"github.com/charmbracelet/lipgloss"
)

// var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

var (
	renderSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render
	renderError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
	renderGray    = lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render
)

type TUI struct {
	ctx      context.Context
	eventsCh chan Event
	loader   *tstatus.StatusWidget
	auth     *tauth.AuthWidget
	chat     *tpage.PageWidget
}

func New(ctx context.Context, eventsCh chan Event) *TUI {
	return &TUI{
		ctx:      ctx,
		eventsCh: eventsCh,
		loader:   tstatus.New(),
		auth:     tauth.New(),
		chat:     tpage.NewPageWidget(),
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

// func (t *TUI) RenderLoader() {
// 	t.loader.SetText("")
// 	t.loader.Run()
// }

func (t *TUI) RenderChat() {
	// t.chat:= tchat.NewPageWidget()
	t.chat.Run()
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
				// t.loader.SetSpinner(event.text)
				t.chat.AddMessage(event.text)
			case eventTypeBar:
				// t.loader.SetProgress(event.text, event.percent)
				msg := fmt.Sprintf("Loading %s: %d%%", event.text, int(event.percent))
				t.chat.AddMessage(msg)
			case eventTypeText:
				// t.loader.SetText(event.text)
				t.chat.AddMessage(event.text)
			case eventTypeTextError:
				// t.loader.SetText(event.text)
				t.chat.AddMessage(renderError(event.text))
			case eventTypeAccess:
				// t.loader.SetAccess(event.access.key, event.access.password)
				output := "\n" + "Client connection credentials:"
				output += "\n" + renderGray("=======================================")
				output += "\n" + "Key\n" + renderSuccess(event.access.key)
				output += "\n\n" + "Password\n" + renderSuccess(event.access.password)
				output += "\n" + renderGray("=======================================")
				t.chat.AddMessage(output)
			}
		}
	}
}
