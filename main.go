package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("S-Conversion")

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to S-Conversion!"),
		widget.NewButton("Single Image", func() {
			print("Single Image")
			// Logic for single image conversion
		}),
		widget.NewButton("Batch Image", func() {
			print("Batch Image")
			// Logic for batch image conversion
		}),
		widget.NewButton("Start Conversion", func() {
			print("Start Conversion")
			// Logic to start conversion
		}),
		widget.NewButton("Close", func() {
			myWindow.Close()
		}),
	))

	myWindow.ShowAndRun()
}
