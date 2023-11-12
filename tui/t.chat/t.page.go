package tchat

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
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

//===== Page

type PageMode int

const (
	ModeLoading PageMode = iota
	ModeAccess
	ModeChat
	ModeMenu
)

type PageWidget struct {
	mode     PageMode
	ready    bool
	messages []*ChatMessage

	topBar   viewport.Model
	viewport viewport.Model
	textarea textarea.Model
	status   viewport.Model

	// menu
	menu     list.Model
	choice   string
	quitting bool
}

func NewPageWidget() *PageWidget {
	ta := textarea.New()
	ta.Placeholder = "Enter text..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(4)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.KeyMap.InsertNewline.SetEnabled(false)

	ta.ShowLineNumbers = false

	messages := make([]*ChatMessage, 0)
	// fake
	// for i := 1; i <= 10; i++ {
	i := rand.Intn(100)
	msg := NewChatMessage("name", fmt.Sprintf("message %d", i))
	msg.SetSent()
	messages = append(messages, msg)
	go func(msg *ChatMessage) {
		time.Sleep(1 * time.Second)
		msg.SetSuccess()
	}(msg)

	i = rand.Intn(100)
	msg = NewChatMessage("", fmt.Sprintf("message %d", i))
	msg.SetSent()
	messages = append(messages, msg)
	go func(msg *ChatMessage) {
		time.Sleep(1 * time.Second)
		msg.SetError()
	}(msg)

	i = rand.Intn(100)
	msg = NewChatMessage("name", fmt.Sprintf("message %d", i))
	msg.SetSent()
	messages = append(messages, msg)
	go func(msg *ChatMessage) {
		time.Sleep(1 * time.Second)
		msg.SetSuccess()
	}(msg)

	i = rand.Intn(100)
	msg = NewChatMessage("", fmt.Sprintf("message %d", i))
	msg.SetSent()
	messages = append(messages, msg)
	go func(msg *ChatMessage) {
		time.Sleep(1 * time.Second)
		msg.SetError()
	}(msg)

	// init menu
	items := []list.Item{
		item("Yes"),
		item("Cancel"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Quit?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	w := PageWidget{
		// mode:     Loading,
		mode:     ModeChat,
		textarea: ta,
		messages: messages,
		menu:     l,
	}
	return &w
}

func (m *PageWidget) Content() string {
	messages := make([]string, 0)
	for _, msg := range m.messages {
		messages = append(messages, msg.Text())
	}
	return strings.Join(messages, "\n\n")
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

func (m *PageWidget) SetMode(mode PageMode) {
	m.mode = mode
}

func (m *PageWidget) AddMessage(text string) {
	// fmt.Printf("len messages: %d\n", len(m.messages))
	msg := NewChatMessage("", text)
	msg.SetSent()
	m.messages = append(m.messages, msg)
	// fmt.Printf("len messages: %d\n", len(m.messages))
}

func (m *PageWidget) Init() tea.Cmd {
	return textarea.Blink
}

func (m *PageWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEsc:
			if m.mode == ModeMenu {
				m.mode = ModeChat
			} else {
				m.mode = ModeMenu
			}

		case tea.KeyCtrlC:
			// fmt.Println("EXIT...")
			return m, tea.Quit

		case tea.KeyTab:
			m.textarea.SetValue(m.textarea.Value() + "\n")

		case tea.KeyEnter:
			// m.messages = append(m.messages, m.textarea.Value())
			m.AddMessage(m.textarea.Value())
			m.viewport.SetContent(m.Content())
			m.textarea.Reset()
			m.viewport.GotoBottom()

		case tea.KeyUp:
			m.menu.CursorUp()

		case tea.KeyDown:
			m.menu.CursorDown()
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
			status := styleOnline.Render("online")
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

	// input
	var cmdText tea.Cmd

	if m.mode == ModeChat {
		// Handle keyboard and mouse events in the viewport
		m.viewport, cmd = m.viewport.Update(msg)

		m.textarea, cmdText = m.textarea.Update(msg)
	}
	var cmdMenu tea.Cmd
	// if m.mode == ModeMenu {
	// 	_, cmd = m.menu.Update(msg)
	// }

	// status
	var cmdStatus tea.Cmd
	m.status, cmdStatus = m.status.Update(msg)
	cmds = append(cmds, cmd, cmdText, cmdStatus, cmdMenu)

	return m, tea.Batch(cmds...)
}

func (m *PageWidget) View() string {
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

	switch m.mode {
	case ModeChat:
		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", m.topBar.View(), m.viewport.View(), m.textarea.View(), m.status.View())

	// MENU
	case ModeMenu:
		if m.choice != "" {
			return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
		}
		if m.quitting {
			return quitTextStyle.Render("Not hungry? That’s cool.")
		}
		return "\n" + m.menu.View()

	// case Loading:
	// return fmt.Sprintf("%s\n%s\n\n%s", m.topBar.View(), m.viewport.View(), m.status.View())
	default:
		return ""
	}
}
