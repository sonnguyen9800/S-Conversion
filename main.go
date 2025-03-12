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

	// Create window
	myWindow := myApp.NewWindow("S-Conversion")

	// Set window icon (this will replace the NO DC text)
	myWindow.SetIcon(theme.FileImageIcon())

	// Set window properties
	myWindow.Resize(fyne.NewSize(800, 600))
	// myWindow.SetFixedSize(false)
	// myWindow.CenterOnScreen()
	// myWindow.SetPadded(false) // Remove default padding

	// Create content with its own padding
	content := view.CreateUI()

	// Set content
	myWindow.SetContent(content)

	// Show and run
	myWindow.ShowAndRun()
}
