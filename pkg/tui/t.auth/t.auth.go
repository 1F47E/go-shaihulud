package tauth

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	enterKey step = iota
	enterPass
	submit
)

const (
	authKeyLen  = 184
	authPassLen = 9
)

var (
	renderSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render
	renderError   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
)

type step int

// type Auth struct {
// 	Key      string
// 	Password string
// }

// func (a Auth) String() string {
// 	return fmt.Sprintf("Key: %s\nPassword: %s", a.Key, a.Password)
// }

type AuthWidget struct {
	viewport viewport.Model
	step     step
	// auth     Auth
	key      string
	password string
	textarea textarea.Model
	value    string
}

func New() *AuthWidget {
	ta := textarea.New()
	ta.Placeholder = "Enter access key..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(4)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(50, 4)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return &AuthWidget{
		textarea: ta,
		step:     enterKey,
		// auth:     Auth{},
		viewport: vp,
	}
}

func (w *AuthWidget) Init() tea.Cmd {
	return textarea.Blink
}

func (w *AuthWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// quit after submit
	if w.step == submit {
		return w, tea.Quit
	}

	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	w.textarea, tiCmd = w.textarea.Update(msg)
	w.viewport, vpCmd = w.viewport.Update(msg)

	// cleanup the input
	text := w.textarea.Value()
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", "")
	w.value = text

	output := ""
	switch w.step {
	case enterKey:
		w.textarea.Placeholder = "Enter access key..."
		if len(text) == authKeyLen {
			output = renderSuccess("Key valid, press enter")
		} else if len(text) > 0 {
			output = renderError(fmt.Sprintf("Access key... %d/%d", len(w.value), authKeyLen))
		}
	case enterPass:
		w.textarea.Placeholder = "Enter password..."
		w.textarea.SetHeight(1)
		if len(text) == authPassLen {
			output = renderSuccess("Password looks valid, press enter")
		} else if len(text) > 0 {
			output = renderError(fmt.Sprintf("Password... %d/%d", len(w.value), authPassLen))
		}
	}

	w.viewport.SetContent(output)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(w.textarea.Value())
			return w, tea.Quit
		case tea.KeyEnter:

			switch w.step {
			case enterKey:
				if len(text) == authKeyLen {
					w.textarea.Reset()
					w.viewport.GotoBottom()
					w.key = text
					w.step = enterPass
				}

			case enterPass:
				if len(text) == authPassLen {
					w.textarea.Reset()
					w.viewport.GotoBottom()
					w.password = text
					w.step = submit
				}
			}
		}
	}

	return w, tea.Batch(tiCmd, vpCmd)
}

func (w *AuthWidget) View() string {
	if w.step == submit {
		return ""
	}
	return fmt.Sprintf(
		"%s\n\n%s",
		w.viewport.View(),
		w.textarea.View(),
	) + "\n\n"
}

func (w *AuthWidget) Run() (string, string, error) {
	p := tea.NewProgram(w)

	if _, err := p.Run(); err != nil {
		return "", "", err
	}
	return w.key, w.password, nil
}
