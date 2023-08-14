package tchat

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

// type ChatWidget struct {
// 	viewport viewport.Model
// 	step     step
// 	messages []string
// 	// auth     Auth
// 	// key      string
// 	// password string
// 	textarea textarea.Model
// 	value    string
// }

// func NewChatWidget() *ChatWidget {

var senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
var styleOnline = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render

type PageWidget struct {
	content  string
	ready    bool
	messages []string

	topBar   viewport.Model
	viewport viewport.Model
	textarea textarea.Model
	status   viewport.Model
	// textarea textarea.Model
	// value    string
}

func NewPageWidget() *PageWidget {
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

	messages := make([]string, 0)
	for i := 1; i <= 100; i++ {
		messages = append(messages, senderStyle.Render("You: ")+fmt.Sprintf("message %d", i))
	}

	content := strings.Join(messages, "\n\n")

	w := PageWidget{
		content:  content,
		textarea: ta,
		messages: messages,
	}
	return &w
}

func (m *PageWidget) Run() {
	// Load some text for our viewport
	// content, err := os.ReadFile("README.md")
	// if err != nil {
	// 	fmt.Println("could not load file:", err)
	// 	os.Exit(1)
	// }

	// fake messages

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

func (m PageWidget) Init() tea.Cmd {
	return textarea.Blink
}

func (m PageWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// headerHeight := lipgloss.Height(m.headerView())
		// footerHeight := lipgloss.Height(m.status.View())
		headerHeight := 2
		footerHeight := 2
		chatHeight := 5 + 1 + 1 // with paddings
		verticalMarginHeight := headerHeight + footerHeight + chatHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// top bar
			m.topBar = viewport.New(msg.Width, 1)
			m.topBar.YPosition = 0
			m.topBar.HighPerformanceRendering = useHighPerformanceRenderer
			m.topBar.SetContent("Chat")

			// chat
			m.textarea.Placeholder = "Say hello..."

			// status
			m.status = viewport.New(msg.Width, 1)
			m.status.YPosition = msg.Height - footerHeight
			m.status.HighPerformanceRendering = useHighPerformanceRenderer
			status := styleOnline("online")
			m.status.SetContent("status: " + status)

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)

	// input
	var cmdText tea.Cmd
	m.textarea, cmdText = m.textarea.Update(msg)

	// status
	var cmdStatus tea.Cmd
	m.status, cmdStatus = m.status.Update(msg)
	cmds = append(cmds, cmd, cmdText, cmdStatus)

	return m, tea.Batch(cmds...)
}

func (m PageWidget) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	// m.viewport.Style = lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.RoundedBorder()).
	// 	BorderForeground(lipgloss.Color("62")).
	// 	PaddingRight(2)

	m.topBar.Style = lipgloss.NewStyle().
		// BorderStyle(lipgloss.RoundedBorder()).
		// BorderForeground(lipgloss.Color("33")).
		Background(lipgloss.Color("99")).
		PaddingRight(2)

	m.status.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#555")).
		// Background(lipgloss.Color("3")).
		PaddingRight(2)

	// return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", m.topBar.View(), m.viewport.View(), m.textarea.View(), m.status.View())
}

// func (m *PageWidget) headerView() string {
// 	title := titleStyle.Render("CHAT")
// 	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
// 	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
// }

// func (m *PageWidget) footerView() string {
// 	// info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
// 	info := "connected"
// 	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
// 	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
// }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func main() {
// 	// Load some text for our viewport
// 	content, err := os.ReadFile("artichoke.md")
// 	if err != nil {
// 		fmt.Println("could not load file:", err)
// 		os.Exit(1)
// 	}

// 	w := PageWidget{content: string(content)}
// 	p := tea.NewProgram(
// 		w,
// 		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
// 		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
// 	)

// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("could not run program:", err)
// 		os.Exit(1)
// 	}
// }
