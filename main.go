package main

import (
	"s_conversion/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	myApp := app.NewWithID("com.s-conversion.app")
	
	// Set app icon
	myApp.SetIcon(theme.FileImageIcon())
	
	myWindow := myApp.NewWindow("S-Conversion")
	
	// Set window icon (this will replace the NO DC text)
	myWindow.SetIcon(theme.FileImageIcon())
	
	// Set window size
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()

	// Set content
	myWindow.SetContent(view.CreateUI())
	
	// Show and run
	myWindow.ShowAndRun()
}
