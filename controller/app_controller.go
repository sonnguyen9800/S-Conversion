package controller

import (
	"fmt"
	"os/exec"
	"runtime"
	"s_conversion/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type AppController struct {
	converter    *model.Converter
	window      fyne.Window
	onProgress  func(float64)
	isConverting bool
}

func NewAppController(window fyne.Window) *AppController {
	return &AppController{
		converter: model.NewConverter(),
		window:   window,
	}
}

func (c *AppController) SetProgressCallback(callback func(float64)) {
	c.onProgress = callback
}

func (c *AppController) HandleSingleImageSelection() (string, error) {
	if c.isConverting {
		return "", fmt.Errorf("conversion in progress")
	}

	// Use PowerShell for Windows, osascript for macOS, or zenity for Linux
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe", "-Command", `Add-Type -AssemblyName System.Windows.Forms
		$f = New-Object System.Windows.Forms.OpenFileDialog
		$f.Filter = "WebP files (*.webp)|*.webp"
		$f.ShowDialog()
		$f.FileName`)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `choose file of type {"webp"} with prompt "Choose a WebP file"`)
	default: // Linux and others
		cmd = exec.Command("zenity", "--file-selection", "--file-filter=*.webp")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error selecting file: %v", err)
	}

	path := string(output)
	if runtime.GOOS == "windows" {
		// Remove newline characters from PowerShell output
		path = path[:len(path)-2]
	}

	c.converter.ConversionType = model.SingleImage
	c.converter.SourcePath = path
	return path, nil
}

func (c *AppController) HandleBatchSelection() (string, error) {
	if c.isConverting {
		return "", fmt.Errorf("conversion in progress")
	}

	// Use PowerShell for Windows, osascript for macOS, or zenity for Linux
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe", "-Command", `Add-Type -AssemblyName System.Windows.Forms
		$f = New-Object System.Windows.Forms.FolderBrowserDialog
		$f.ShowDialog()
		$f.SelectedPath`)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `choose folder with prompt "Choose a folder containing WebP files"`)
	default: // Linux and others
		cmd = exec.Command("zenity", "--file-selection", "--directory")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error selecting folder: %v", err)
	}

	path := string(output)
	if runtime.GOOS == "windows" {
		// Remove newline characters from PowerShell output
		path = path[:len(path)-2]
	}

	c.converter.ConversionType = model.BatchImage
	c.converter.SourcePath = path
	return path, nil
}

func (c *AppController) HandleOutputSelection() (string, error) {
	if c.isConverting {
		return "", fmt.Errorf("conversion in progress")
	}

	// Use PowerShell for Windows, osascript for macOS, or zenity for Linux
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell.exe", "-Command", `Add-Type -AssemblyName System.Windows.Forms
		$f = New-Object System.Windows.Forms.FolderBrowserDialog
		$f.Description = "Select output folder for converted images"
		$f.ShowDialog()
		$f.SelectedPath`)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `choose folder with prompt "Select output folder for converted images"`)
	default: // Linux and others
		cmd = exec.Command("zenity", "--file-selection", "--directory", "--title=Select output folder")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error selecting output folder: %v", err)
	}

	path := string(output)
	if runtime.GOOS == "windows" {
		// Remove newline characters from PowerShell output
		path = path[:len(path)-2]
	}

	c.converter.OutputPath = path
	return path, nil
}

func (c *AppController) StartConversion() error {
	if c.isConverting {
		return fmt.Errorf("conversion already in progress")
	}

	if c.converter.SourcePath == "" {
		return fmt.Errorf("no source selected")
	}

	c.isConverting = true
	var err error

	go func() {
		defer func() {
			c.isConverting = false
			if c.onProgress != nil {
				c.onProgress(1.0)
			}
		}()

		switch c.converter.ConversionType {
		case model.SingleImage:
			err = c.converter.ConvertSingle(c.converter.SourcePath)
		case model.BatchImage:
			err = c.converter.ConvertBatch(c.converter.SourcePath)
		default:
			err = fmt.Errorf("invalid conversion type")
		}

		if err != nil {
			dialog.ShowError(err, c.window)
			return
		}

		dialog.ShowInformation("Success", "Conversion completed successfully!", c.window)
	}()

	return nil
}

func (c *AppController) IsConverting() bool {
	return c.isConverting
}

func (c *AppController) GetProgress() float64 {
	return c.converter.Progress
}
