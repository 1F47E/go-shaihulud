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

var senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
var styleOnline = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render

type PageWidget struct {
	ready    bool
	messages []string

	topBar   viewport.Model
	viewport viewport.Model
	textarea textarea.Model
	status   viewport.Model
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
	ta.KeyMap.InsertNewline.SetEnabled(false)

	ta.ShowLineNumbers = false

	messages := make([]string, 0)
	for i := 1; i <= 100; i++ {
		messages = append(messages, senderStyle.Render("You: ")+fmt.Sprintf("message %d", i))
	}

	w := PageWidget{
		textarea: ta,
		messages: messages,
	}
	return &w
}

func (m *PageWidget) Content() string {
	return strings.Join(m.messages, "\n\n")
}

func (m *PageWidget) Run() {
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // fullscreen
		tea.WithMouseCellMotion(), // mouse scroll for messages
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

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println("EXIT...")
			return m, tea.Quit

		case tea.KeyTab:
			m.textarea.SetValue(m.textarea.Value() + "\n")

		case tea.KeyEnter:
			m.messages = append(m.messages, m.textarea.Value())
			m.viewport.SetContent(m.Content())
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
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
			m.viewport.SetContent(m.Content())
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
	default:
		// update viewport with messages
		m.viewport.SetContent(m.Content())
		m.viewport.GotoBottom()
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

	m.topBar.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("99")).
		PaddingRight(2)

	m.status.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#555")).
		PaddingRight(2)

	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", m.topBar.View(), m.viewport.View(), m.textarea.View(), m.status.View())
}
