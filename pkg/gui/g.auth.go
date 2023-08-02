package gui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Form(nextSlide func()) (title string, content tview.Primitive) {
	form := tview.NewForm().
		AddTextArea("ACCESS KEY:", "", 0, 6, 0, nil).
		AddPasswordField("PASSWORD", "", 10, '*', nil).
		AddButton("Connect", func() {
			fmt.Println("Connect")
		})

	// chatBox := tview.NewBox().
	// 	SetTitle("Chat window").
	// 	SetBorder(true).
	// 	SetBorderAttributes(tcell.AttrDim)
	width := 50
	height := 31
	chatWIndowFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		// AddItem(form, 0, 6, true)
		AddItem(Center(width, height, form), 0, 1, true)
	return "form", chatWIndowFlex
}

func Auth(nextSlide func()) (title string, content tview.Primitive) {
	b, _ := base64.StdEncoding.DecodeString(IMG_AUTH)
	img, _ := jpeg.Decode(bytes.NewReader(b))
	image := tview.NewImage().
		SetImage(img)
	// image.SetColors(256)
	image.SetColors(tview.TrueColor)
	image.SetRect(0, 0, 24, 24)

	form := tview.NewForm().
		// AddImage("Photo:", img, 0, 12, 0).
		AddTextArea("Access key:", "", 0, 6, 0, nil).
		AddPasswordField("Password:", "", 10, '*', nil).
		AddButton("Connect", nextSlide)
	form.SetFocus(0)
	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorGray)
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetButtonTextColor(tcell.ColorBlack)
	form.SetButtonBackgroundColor(tcell.ColorGreen)
	button := tview.NewButton("Hit Enter to close").SetSelectedFunc(func() {
		app.Stop()
	})
	button.SetBorder(true).SetRect(0, 0, 22, 3)

	width := 50
	height := 31
	window := tview.NewFlex().
		// AddItem(image, 0, 1, true).
		AddItem(Center(width, height, form), 0, 1, true)

	return "auth", window

}
func Center(width, height int, p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
