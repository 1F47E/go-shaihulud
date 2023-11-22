package tstatus

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding        = 2
	maxWidth       = 80
	colorPrimary   = "#ABCD03"
	colorSecondary = "#7D4698"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

type mode int

const (
	spin mode = iota
	bar
	text
)

type StatusWidget struct {
	mode         mode
	title        string
	authKey      string
	authPassword string
	spinner      spinner.Model
	progress     progress.Model
	percent      float64
}

func New() *StatusWidget {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))

	return &StatusWidget{
		mode:     text,
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
		percent:  0,
	}
}

func (w *StatusWidget) SetProgress(title string, percent float64) {
	w.mode = bar
	w.title = title
	w.percent = percent
}

func (w *StatusWidget) SetSpinner(title string) {
	w.mode = spin
	w.title = title
}

func (w *StatusWidget) SetText(title string) {
	w.mode = text
	w.title = title
}

func (w *StatusWidget) SetAccess(authKey, authPassword string) {
	w.authKey = authKey
	w.authPassword = authPassword
}

func (w *StatusWidget) Run() {
	isFullScreen := false
	var p *tea.Program
	if isFullScreen {
		p = tea.NewProgram(w, tea.WithAltScreen())
	} else {
		p = tea.NewProgram(w)
	}

	if _, err := p.Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func (w *StatusWidget) Init() tea.Cmd {
	return tea.Batch(tickCmd(), w.spinner.Tick)
}

func (w *StatusWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TOOD: remove?
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			return w, tea.Quit

		default:
			var cmd tea.Cmd
			return w, cmd
		}

	case tea.WindowSizeMsg:
		w.progress.Width = msg.Width - padding*2 - 4
		if w.progress.Width > maxWidth {
			w.progress.Width = maxWidth
		}
		return w, nil

	case tickMsg:
		cmd := w.progress.SetPercent(w.percent)
		return w, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := w.progress.Update(msg)
		w.progress = progressModel.(progress.Model)
		return w, cmd

	default:
		var cmd tea.Cmd
		w.spinner, cmd = w.spinner.Update(msg)
		return w, cmd
	}
}

func (w *StatusWidget) View() string {

	pad := strings.Repeat(" ", padding)

	output := ""
	if w.authKey != "" && w.authPassword != "" {
		output += "\n" + "Client connection credentials:"
		output += "\n" + helpStyle("=======================================")
		output += "\n" + "Key\n" + w.authKey
		output += "\n\n" + "Password\n" + w.authPassword
		output += "\n" + helpStyle("=======================================")
		output += "\n"
	}

	if w.mode == text {
		output += fmt.Sprintf("\n\n%s%s\n\n", pad, w.title)
	} else if w.mode == spin {
		output += fmt.Sprintf("\n\n%s%s %s\n\n", pad, w.spinner.View(), w.title)

	} else if w.mode == bar {
		output += "\n" +
			pad + w.title + "\n\n" +
			pad + w.progress.View() + "\n"
	}
	return output
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Microsecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
