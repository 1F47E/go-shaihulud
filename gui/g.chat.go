package gui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Submit() {

}

// Flex demonstrates flexbox layout.
func Chat(nextSlide func()) (title string, content tview.Primitive) {
	// modalShown := false
	pages := tview.NewPages()

	logsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetTextColor(tcell.ColorYellow).
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	logsView.SetBorder(true)
	logsView.SetTitle("Logs")

	fmt.Fprint(logsView, "[red]Some logs[white]\n")

	// chat window row
	chatBox := tview.NewBox().
		SetTitle("Chat window").
		SetBorder(true).
		SetBorderAttributes(tcell.AttrDim)

	inputForm := tview.NewForm().
		AddInputField("Say:", "", 80, nil, nil).
		// SetBackgroundColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetFocus(0)

	// catch enter press
	inputForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCR {
			fmt.Fprintf(logsView, "[yellow]ENTER PRESSED[white]\n")
			fmt.Fprintf(logsView, "data: %v\n", inputForm.GetFormItem(0).(*tview.InputField).GetText())
			// clear input
			inputForm.GetFormItem(0).(*tview.InputField).SetText("")
			return nil
		}
		return event
	})

	inputBox := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(inputForm, 0, 12, true)

	chatWindowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatBox, 0, 8, false).
		AddItem(inputBox, 0, 1, true)

	flex := tview.NewFlex().
		AddItem(chatWindowFlex, 0, 2, true).
		// SetDirection(tview.FlexRow).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle("Input"), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			// AddItem(tview.NewBox().SetBorder(true).SetTitle("Flexible width"), 0, 1, false).
			AddItem(logsView, 0, 1, false), 0, 1, false)
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
