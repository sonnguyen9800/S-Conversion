package controller

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"s_conversion/model"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type AppController struct {
	converter          *model.Converter
	window             fyne.Window
	onProgress         func(float64)
	onStatus           func(string)
	onOutputPathChange func(string)
	isConverting       bool
}

func NewAppController(window fyne.Window) *AppController {
	return &AppController{
		converter: model.NewConverter(),
		window:    window,
	}
}

func (c *AppController) SetProgressCallback(callback func(float64)) {
	c.onProgress = callback
}

func (c *AppController) SetStatusCallback(callback func(string)) {
	c.onStatus = callback
}

func (c *AppController) SetOutputPathCallback(callback func(string)) {
	c.onOutputPathChange = callback
}

func (c *AppController) updateStatus(status string) {
	if c.onStatus != nil {
		c.onStatus(status)
	}
}

func (c *AppController) updateOutputPath() {
	if c.onOutputPathChange == nil {
		return
	}

	outputPath := c.converter.OutputPath
	if outputPath == "" {
		if c.converter.SourcePath == "" {
			c.onOutputPathChange("")
			return
		}

		if c.converter.ConversionType == model.SingleImage {
			outputPath = filepath.Dir(c.converter.SourcePath)
		} else {
			outputPath = c.converter.SourcePath
		}
	}
	c.onOutputPathChange(outputPath)
}

func (c *AppController) OpenOutputFolder() error {
	outputPath := c.converter.OutputPath
	if outputPath == "" {
		if c.converter.SourcePath == "" {
			return fmt.Errorf("no output location available")
		}

		if c.converter.ConversionType == model.SingleImage {
			outputPath = filepath.Dir(c.converter.SourcePath)
		} else {
			outputPath = c.converter.SourcePath
		}
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("output folder does not exist")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", outputPath)
	case "darwin":
		cmd = exec.Command("open", outputPath)
	default: // Linux and others
		cmd = exec.Command("xdg-open", outputPath)
	}
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// Ignore exit status 1 since Explorer opened successfully
			return nil
		}
		fmt.Printf("error opening folder: %v\n", err)
		return fmt.Errorf("error opening folder: %v", err)
	}
	return nil
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
		if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
			$f.FileName
		}`)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `choose file of type {"webp"} with prompt "Choose a WebP file"`)
	default: // Linux and others
		cmd = exec.Command("zenity", "--file-selection", "--file-filter=*.webp")
	}

	output, err := cmd.Output()
	if err != nil {
		// Check if it's just a cancel operation
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			c.updateStatus("No file selected")
			return "", fmt.Errorf("file selection cancelled")
		}
		c.updateStatus("Error selecting file")
		return "", fmt.Errorf("error selecting file: %v", err)
	}

	path := string(output)
	path = strings.TrimSpace(path) // Remove any whitespace, newlines

	// Handle platform-specific path formatting
	switch runtime.GOOS {
	case "windows":
		// Remove any "OK" or dialog result text that might appear
		if parts := strings.Split(path, "\r\n"); len(parts) > 0 {
			path = parts[len(parts)-1] // Take the last non-empty line
		}
	case "darwin":
		// macOS osascript might return alias format, convert to regular path
		path = strings.TrimPrefix(path, "alias ")
		path = strings.Trim(path, "\n")
	}

	if path == "" {
		c.updateStatus("No file selected")
		return "", fmt.Errorf("no file selected")
	}

	c.converter.ConversionType = model.SingleImage
	c.converter.SourcePath = path
	c.updateStatus(fmt.Sprintf("Selected file: %s", filepath.Base(path)))
	c.updateOutputPath()
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
		if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
			$f.SelectedPath
		}`)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `choose folder with prompt "Choose a folder containing WebP files"`)
	default: // Linux and others
		cmd = exec.Command("zenity", "--file-selection", "--directory")
	}

	output, err := cmd.Output()
	if err != nil {
		// Check if it's just a cancel operation
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			c.updateStatus("No folder selected")
			return "", fmt.Errorf("folder selection cancelled")
		}
		c.updateStatus("Error selecting folder")
		return "", fmt.Errorf("error selecting folder: %v", err)
	}

	path := string(output)
	path = strings.TrimSpace(path) // Remove any whitespace, newlines

	// Handle platform-specific path formatting
	switch runtime.GOOS {
	case "windows":
		// Remove any "OK" or dialog result text that might appear
		if parts := strings.Split(path, "\r\n"); len(parts) > 0 {
			path = parts[len(parts)-1] // Take the last non-empty line
		}
	case "darwin":
		// macOS osascript might return alias format, convert to regular path
		path = strings.TrimPrefix(path, "alias ")
		path = strings.Trim(path, "\n")
	}

	if path == "" {
		c.updateStatus("No folder selected")
		return "", fmt.Errorf("no folder selected")
	}

	c.converter.ConversionType = model.BatchImage
	c.converter.SourcePath = path
	c.updateStatus(fmt.Sprintf("Selected folder: %s", filepath.Base(path)))
	c.updateOutputPath()
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
		if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
			$f.SelectedPath
		}`)
	case "darwin":
		cmd = exec.Command("osascript", "-e", `choose folder with prompt "Select output folder for converted images"`)
	default: // Linux and others
		cmd = exec.Command("zenity", "--file-selection", "--directory", "--title=Select output folder")
	}

	output, err := cmd.Output()
	if err != nil {
		// Check if it's just a cancel operation
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", fmt.Errorf("output folder selection cancelled")
		}
		return "", fmt.Errorf("error selecting output folder: %v", err)
	}

	path := string(output)
	path = strings.TrimSpace(path) // Remove any whitespace, newlines

	// Handle platform-specific path formatting
	switch runtime.GOOS {
	case "windows":
		// Remove any "OK" or dialog result text that might appear
		if parts := strings.Split(path, "\r\n"); len(parts) > 0 {
			path = parts[len(parts)-1] // Take the last non-empty line
		}
	case "darwin":
		// macOS osascript might return alias format, convert to regular path
		path = strings.TrimPrefix(path, "alias ")
		path = strings.Trim(path, "\n")
	}

	if path == "" {
		return "", fmt.Errorf("no output folder selected")
	}

	c.converter.OutputPath = path
	c.updateOutputPath()
	return path, nil
}

func (c *AppController) StartConversion() error {
	if c.isConverting {
		return fmt.Errorf("conversion already in progress")
	}

	if c.converter.SourcePath == "" {
		c.updateStatus("No source selected")
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
			c.updateStatus(fmt.Sprintf("Converting %s...", filepath.Base(c.converter.SourcePath)))
			err = c.converter.ConvertSingle(c.converter.SourcePath)
		case model.BatchImage:
			c.updateStatus(fmt.Sprintf("Converting files in folder: %s...", filepath.Base(c.converter.SourcePath)))
			err = c.converter.ConvertBatch(c.converter.SourcePath)
		default:
			err = fmt.Errorf("invalid conversion type")
		}

		if err != nil {
			dialog.ShowError(err, c.window)
			c.updateStatus("Conversion failed")
			return
		}

		// Show success message with output location
		outputLocation := c.converter.OutputPath
		if outputLocation == "" {
			if c.converter.ConversionType == model.SingleImage {
				outputLocation = filepath.Dir(c.converter.SourcePath)
			} else {
				outputLocation = c.converter.SourcePath
			}
		}

		successMsg := fmt.Sprintf("Conversion completed successfully!")
		dialog.ShowInformation("Success", successMsg, c.window)
		c.updateStatus("Conversion completed")
		c.updateOutputPath()
	}()

	return nil
}

func (c *AppController) IsConverting() bool {
	return c.isConverting
}

func (c *AppController) GetProgress() float64 {
	return c.converter.Progress
}
