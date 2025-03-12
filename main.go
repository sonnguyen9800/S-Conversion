package main

import (
	"s_conversion/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("S-Conversion")

	myWindow.SetContent(view.CreateUI())
	myWindow.ShowAndRun()
}
