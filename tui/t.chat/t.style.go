package tchat

import "github.com/charmbracelet/lipgloss"

var (
	styleSender   = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	styleReceiver = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	styleRed      = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	styleGreen    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	styleYellow   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleBlue     = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	styleGray     = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	styleWhite    = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	styleBlack    = lipgloss.NewStyle().Foreground(lipgloss.Color("16"))
	styleOnline   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)
