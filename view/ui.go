package view

import (
	"fmt"
	"path/filepath"
	"s_conversion/controller"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	window     fyne.Window
	controller *controller.AppController
	progress   *widget.ProgressBar
	status     *widget.Label
	outputPath *widget.Label
	openFolderBtn *widget.Button
}

func truncatePath(path string) string {
	if len(path) <= 40 {
		return path
	}
	
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	
	// If the base name itself is too long
	if len(base) > 20 {
		base = base[:17] + "..."
	}
	
	// Get the first and last parts of the directory
	parts := filepath.SplitList(dir)
	if len(parts) <= 2 {
		return filepath.Join(dir, base)
	}
	
	return filepath.Join(parts[0], "...", parts[len(parts)-1], base)
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func NewUI(window fyne.Window) *UI {
	ui := &UI{
		window:     window,
		controller: controller.NewAppController(window),
		progress:   widget.NewProgressBar(),
		status:     widget.NewLabel("No file selected"),
		outputPath: widget.NewLabel(""),
	}

	ui.openFolderBtn = widget.NewButtonWithIcon("Open Destination Folder", theme.FolderOpenIcon(), func() {
		if err := ui.controller.OpenOutputFolder(); err != nil {
			dialog.ShowError(err, window)
		}
	})
	ui.openFolderBtn.Hide()

	ui.controller.SetProgressCallback(func(progress float64) {
		ui.progress.SetValue(progress)
	})

	ui.controller.SetStatusCallback(func(status string) {
		if strings.HasPrefix(status, "Selected file: ") {
			filename := strings.TrimPrefix(status, "Selected file: ")
			ui.status.SetText("Selected file: " + truncateText(filename, 30))
		} else if strings.HasPrefix(status, "Selected folder: ") {
			foldername := strings.TrimPrefix(status, "Selected folder: ")
			ui.status.SetText("Selected folder: " + truncateText(foldername, 30))
		} else {
			ui.status.SetText(status)
		}
	})

	ui.controller.SetOutputPathCallback(func(path string) {
		if path == "" {
			ui.outputPath.SetText("")
			ui.openFolderBtn.Hide()
		} else {
			ui.outputPath.SetText(fmt.Sprintf("Output folder: %s", truncatePath(path)))
			ui.openFolderBtn.Show()
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
			dialog.ShowInformation("S-Conversion",
				"S-Conversion is a simple tool to convert WebP images to PNG format made by sonnguyen9800.\n"+
					"Version 1.0.0\n"+
					"© 2024 S-Conversion ",
				u.window)
		}),
	)

	return fyne.NewMainMenu(fileMenu, helpMenu)
}

func (u *UI) handleSingleImageSelection() {
	if u.controller.IsConverting() {
		dialog.ShowError(fmt.Errorf("conversion in progress"), u.window)
		return
	}
	_, err := u.controller.HandleSingleImageSelection()
	if err != nil {
		dialog.ShowError(err, u.window)
	}
}

func (u *UI) handleBatchSelection() {
	if u.controller.IsConverting() {
		dialog.ShowError(fmt.Errorf("conversion in progress"), u.window)
		return
	}
	_, err := u.controller.HandleBatchSelection()
	if err != nil {
		dialog.ShowError(err, u.window)
	}
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
	})

	// Create header with specific styling
	header := widget.NewLabel("Convert WebP to PNG")
	header.TextStyle = fyne.TextStyle{Bold: true}
	
	// Create containers without additional padding
	buttonContainer := container.New(layout.NewVBoxLayout(),
		singleButton,
		batchButton,
	)

	statusContainer := container.New(layout.NewVBoxLayout(),
		container.NewHBox(
			widget.NewLabel("Status:"),
			u.status,
		),
		u.outputPath,
		u.openFolderBtn,
	)

	// Create main content with minimal spacing
	mainContent := container.New(layout.NewVBoxLayout(),
		header,
		layout.NewSpacer(),
		buttonContainer,
		outputButton,
		convertButton,
		u.progress,
		statusContainer,
	)

	// Use a padding container with specific padding values
	return container.NewPadded(mainContent)
}

// CreateUI creates and returns the main UI container
func CreateUI() fyne.CanvasObject {
	window := fyne.CurrentApp().Driver().AllWindows()[0]
	ui := NewUI(window)
	return ui.CreateUI()
}
