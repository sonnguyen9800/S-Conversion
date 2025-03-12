package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateUI creates the user interface components.
func CreateUI() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Welcome to S-Conversion!"),
		widget.NewButton("Single Image", func() {
			// Logic for single image conversion
		}),
		widget.NewButton("Batch Image", func() {
			// Logic for batch image conversion
		}),
		widget.NewButton("Start Conversion", func() {
			// Logic to start conversion
		}),
		widget.NewButton("Close", func() {
			// Logic to close the application
		}),
	)
}
