package model

import (
	"image/png"
	"os"

	"github.com/chai2010/webp"
)

// ConvertWebPToPNG converts a WebP image to PNG format.
func ConvertWebPToPNG(inputPath string, outputPath string) error {
	// Open the WebP image file
	webpFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer webpFile.Close()

	// Decode the WebP image
	img, err := webp.Decode(webpFile)
	if err != nil {
		return err
	}

	// Create a new PNG file
	pngFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer pngFile.Close()

	// Encode the image to PNG format
	if err := png.Encode(pngFile, img); err != nil {
		return err
	}

	return nil
}
