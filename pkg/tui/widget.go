package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// menu stuff
var docStyle = lipgloss.NewStyle().Margin(1, 2)

const listHeight = 14

// loader stuff
const (
	padding  = 2
	maxWidth = 80
	// tor logo colors
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

// menu stuff

// type itemDelegate struct{}

// func (d itemDelegate) Height() int                             { return 1 }
// func (d itemDelegate) Spacing() int                            { return 0 }
// func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
// func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
// 	i, ok := listItem.(item)
// 	if !ok {
// 		return
// 	}

// 	str := fmt.Sprintf("%d. %s", index+1, i)

// 	fn := itemStyle.Render
// 	if index == m.Index() {
// 		fn = func(s ...string) string {
// 			return selectedItemStyle.Render("> " + strings.Join(s, " "))
// 		}
// 	}

// 	fmt.Fprint(w, fn(str))
// }

type Widget struct {
	list   list.Model
	choice string

	mode         mode
	title        string
	authKey      string
	authPassword string
	spinner      spinner.Model
	progress     progress.Model
	percent      float64
}

func NewWidget() *Widget {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimary))

	// MENU
	// items := []list.Item{
	// 	item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
	// 	item{title: "Nutella", desc: "It's good on toast"},
	// 	item{title: "Bitter melon", desc: "It cools you down"},
	// 	item{title: "Nice socks", desc: "And by that I mean socks without holes"},
	// }

	// l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	// l.Title = "My Fave Things"

	return &Widget{
		mode: text,
		// list:     l,
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
		percent:  0,
	}
}

func (w *Widget) SetProgress(title string, percent float64) {
	w.mode = bar
	w.title = title
	w.percent = percent
}

func (w *Widget) SetSpinner(title string) {
	w.mode = spin
	w.title = title
}

func (w *Widget) SetText(title string) {
	w.mode = text
	w.title = title
}

func (w *Widget) SetAccess(authKey, authPassword string) {
	w.authKey = authKey
	w.authPassword = authPassword
}

func (w *Widget) Run() {
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

func (w *Widget) Init() tea.Cmd {
	return tea.Batch(tickCmd(), w.spinner.Tick)
}

// func (w *Widget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if msg.String() == "ctrl+c" {
// 			return w, tea.Quit
// 		}
// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		w.list.SetSize(msg.Width-h, msg.Height-v)
// 	}

// 	var cmd tea.Cmd
// 	w.list, cmd = w.list.Update(msg)
// 	return w, cmd
// }

func (w *Widget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// return w, tea.Quit
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			return w, tea.Quit
		// case "enter":
		// 	i, ok := w.list.SelectedItem().(item)
		// 	if ok {
		// 		w.choice = string(i)
		// 	}
		// 	return w, tea.Quit
		default:
			var cmd tea.Cmd
			// e.viewport, cmd = e.viewport.Update(msg)
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
		// var cmd2 tea.Cmd
		// w.list, cmd2 = w.list.Update(msg)
		// return w, tea.Batch(cmd1, cmd2)
		return w, cmd
	}
}

func (w *Widget) View() string {
	// render menu only for now
	// return "\n" + w.list.View()

	pad := strings.Repeat(" ", padding)

	output := ""
	if w.authKey != "" && w.authPassword != "" {
		output += "\n" + pad + " Client auth creds for connection"
		output += "\n" + pad + helpStyle("=======================================")
		output += "\n" + pad + " Key: " + w.authKey
		output += "\n" + pad + " Password: " + w.authPassword
		output += "\n" + pad + helpStyle("=======================================")
		output += "\n"
	}

	// doc := strings.Builder{}
	// var renderedTabs []string
	// var style lipgloss.Style
	// border, _, _, _, _ := style.GetBorder()
	// style = style.Border(border)
	// renderedTabs = append(renderedTabs, style.Render("TAB1"))
	// row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	// doc.WriteString(row)
	// doc.WriteString("\n")
	// doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
	// return docStyle.Render(doc.String())

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
