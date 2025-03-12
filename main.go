package main

import (
	"s_conversion/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.s-conversion.app")
	myWindow := myApp.NewWindow("S-Conversion")
	
	// Set window size
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()

	// Set content
	myWindow.SetContent(view.CreateUI())
	
	// Show and run
	myWindow.ShowAndRun()
}
