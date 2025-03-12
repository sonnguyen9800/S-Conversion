package controller

import (
	"fmt"
	"s_conversion/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
	
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, c.window)
			return
		}
		if reader == nil {
			return
		}
		c.converter.ConversionType = model.SingleImage
		c.converter.SourcePath = reader.URI().Path()
	}, c.window)
	
	dialog.SetFilter(storage.NewExtensionFileFilter([]string{".webp"}))
	dialog.Show()
	
	return c.converter.SourcePath, nil
}

func (c *AppController) HandleBatchSelection() (string, error) {
	if c.isConverting {
		return "", fmt.Errorf("conversion in progress")
	}

	dialog := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, c.window)
			return
		}
		if list == nil {
			return
		}
		c.converter.ConversionType = model.BatchImage
		c.converter.SourcePath = list.Path()
	}, c.window)
	
	dialog.Show()
	
	return c.converter.SourcePath, nil
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
