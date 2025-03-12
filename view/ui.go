package view

import (
	"fmt"
	"s_conversion/controller"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	window     fyne.Window
	controller *controller.AppController
	progress   *widget.ProgressBar
	status     *widget.Label
}

func NewUI(window fyne.Window) *UI {
	ui := &UI{
		window:     window,
		controller: controller.NewAppController(window),
		progress:   widget.NewProgressBar(),
		status:     widget.NewLabel("Ready"),
	}

	ui.controller.SetProgressCallback(func(progress float64) {
		ui.progress.SetValue(progress)
		if progress >= 1.0 {
			ui.status.SetText("Ready")
		}
	})

	return ui
}

func (u *UI) createMenuBar() *fyne.MainMenu {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Single Image", func() {
			u.handleSingleImageSelection()
		}),
		fyne.NewMenuItem("Batch Conversion", func() {
			u.handleBatchSelection()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Exit", func() {
			u.window.Close()
		}),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			dialog.ShowInformation("About S-Conversion",
				"S-Conversion is a simple tool to convert WebP images to PNG format.\n"+
					"Version 1.0.0\n"+
					"© 2024 S-Conversion",
				u.window)
		}),
	)

	return fyne.NewMainMenu(fileMenu, helpMenu)
}

func (u *UI) handleSingleImageSelection() {
	if u.controller.IsConverting() {
		dialog.ShowError(nil, u.window)
		return
	}
	_, err := u.controller.HandleSingleImageSelection()
	if err != nil {
		dialog.ShowError(err, u.window)
	}
	u.status.SetText("Single image selected")
}

func (u *UI) handleBatchSelection() {
	if u.controller.IsConverting() {
		dialog.ShowError(nil, u.window)
		return
	}
	_, err := u.controller.HandleBatchSelection()
	if err != nil {
		dialog.ShowError(err, u.window)
	}
	u.status.SetText("Folder selected for batch conversion")
}

func (u *UI) handleOutputSelection() {
	if u.controller.IsConverting() {
		dialog.ShowError(fmt.Errorf("conversion in progress"), u.window)
		return
	}
	_, err := u.controller.HandleOutputSelection()
	if err != nil {
		dialog.ShowError(err, u.window)
		return
	}
	u.status.SetText("Output destination selected")
}

func (u *UI) CreateUI() fyne.CanvasObject {
	// Set menu
	u.window.SetMainMenu(u.createMenuBar())

	// Create buttons
	singleButton := widget.NewButtonWithIcon("Single Image", theme.FileIcon(), u.handleSingleImageSelection)
	batchButton := widget.NewButtonWithIcon("Batch Conversion", theme.FolderIcon(), u.handleBatchSelection)
	outputButton := widget.NewButtonWithIcon("Set Output Folder", theme.FolderOpenIcon(), u.handleOutputSelection)
	
	convertButton := widget.NewButtonWithIcon("Start Conversion", theme.MediaPlayIcon(), func() {
		if err := u.controller.StartConversion(); err != nil {
			dialog.ShowError(err, u.window)
			return
		}
		u.status.SetText("Converting...")
	})

	// Create main content
	content := container.NewVBox(
		widget.NewLabel("Welcome to S-Conversion"),
		container.NewHBox(
			singleButton,
			batchButton,
		),
		outputButton,
		convertButton,
		u.progress,
		container.NewHBox(
			widget.NewLabel("Status:"),
			u.status,
		),
	)

	return container.NewPadded(content)
}

// CreateUI creates and returns the main UI container
func CreateUI() fyne.CanvasObject {
	window := fyne.CurrentApp().Driver().AllWindows()[0]
	ui := NewUI(window)
	return ui.CreateUI()
}
