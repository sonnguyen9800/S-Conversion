package controller

import (
	"s_conversion/view"

	"fyne.io/fyne/v2/app"
)

// StartApp initializes the application and
// StartApp initializes the application and
func StartApp() {
	myApp := app.New()
	myWindow := myApp.NewWindow("S-Conversion")

	myWindow.SetContent(view.CreateUI())
	myWindow.ShowAndRun()
}
