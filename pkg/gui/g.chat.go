package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Flex demonstrates flexbox layout.
func Chat(nextSlide func()) (title string, content tview.Primitive) {
	// modalShown := false
	pages := tview.NewPages()

	// chat window row
	chatBox := tview.NewBox().
		SetTitle("Chat window").
		SetBorder(true).
		SetBorderAttributes(tcell.AttrDim)

	chatWIndowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatBox, 0, 6, true)

	// input box
	inputBox := tview.NewBox().
		SetTitle("Input").
		SetBorder(true).
		SetBorderAttributes(tcell.AttrDim)

	chatWIndowFlex.AddItem(inputBox, 0, 1, false)

	// logs
	logsBox := tview.NewBox().
		SetBorder(true).
		SetTitle("Logs")

	flex := tview.NewFlex().
		AddItem(chatWIndowFlex, 0, 2, true).
		// SetDirection(tview.FlexRow).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle("Input"), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			// AddItem(tview.NewBox().SetBorder(true).SetTitle("Flexible width"), 0, 1, false).
			AddItem(logsBox, 0, 1, false), 0, 1, false)
		// AddItem(tview.NewBox().SetBorder(true).SetTitle("Fixed width"), 30, 1, false)

	// modal window
	// flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if modalShown {
	// 		nextSlide()
	// 		modalShown = false
	// 	} else {
	// 		pages.ShowPage("modal")
	// 		modalShown = true
	// 	}
	// 	return event
	// })
	// modal := tview.NewModal().
	// 	SetText("Resize the window to see the effect of the flexbox parameters").
	// 	AddButtons([]string{"Ok"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	// 	pages.HidePage("modal")
	// })
	pages.AddPage("flex", flex, true, true)
	// AddPage("modal", modal, false, false)
	return "Chat", pages
}
