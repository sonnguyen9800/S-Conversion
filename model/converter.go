package model

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"
)

type ConversionType int

const (
	SingleImage ConversionType = iota
	BatchImage
)

const (
	MaxBatchSize    = 100
	MaxTotalSizeMB  = 100
	MaxFileSizeInMB = MaxTotalSizeMB * 1024 * 1024 // 100MB in bytes
)

type Converter struct {
	ConversionType ConversionType
	SourcePath     string
	Progress       float64
	TotalFiles     int
	ConvertedFiles int
}

func NewConverter() *Converter {
	return &Converter{
		Progress: 0,
	}
}

func (c *Converter) ValidateFile(path string) error {
	if !strings.HasSuffix(strings.ToLower(path), ".webp") {
		return errors.New("file is not a WebP image")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error accessing file: %v", err)
	}

	if info.Size() > MaxFileSizeInMB {
		return fmt.Errorf("file size exceeds maximum limit of %d MB", MaxTotalSizeMB)
	}

	return nil
}

func (c *Converter) ConvertSingle(sourcePath string) error {
	if err := c.ValidateFile(sourcePath); err != nil {
		return err
	}

	return c.convertFile(sourcePath)
}

func (c *Converter) ConvertBatch(folderPath string) error {
	files, err := filepath.Glob(filepath.Join(folderPath, "*.webp"))
	if err != nil {
		return fmt.Errorf("error reading folder: %v", err)
	}

	if len(files) > MaxBatchSize {
		return fmt.Errorf("too many files (max %d allowed)", MaxBatchSize)
	}

	var totalSize int64
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("error accessing file %s: %v", file, err)
		}
		totalSize += info.Size()
	}

	if totalSize > MaxFileSizeInMB {
		return fmt.Errorf("total size exceeds maximum limit of %d MB", MaxTotalSizeMB)
	}

	c.TotalFiles = len(files)
	c.ConvertedFiles = 0

	for _, file := range files {
		if err := c.convertFile(file); err != nil {
			return fmt.Errorf("error converting %s: %v", file, err)
		}
		c.ConvertedFiles++
		c.Progress = float64(c.ConvertedFiles) / float64(c.TotalFiles)
	}

	return nil
}

func (c *Converter) convertFile(sourcePath string) error {
	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer sourceFile.Close()

	// Decode WebP
	img, err := webp.Decode(sourceFile)
	if err != nil {
		return fmt.Errorf("error decoding WebP: %v", err)
	}

	// Create output PNG file
	outputPath := strings.TrimSuffix(sourcePath, ".webp") + ".png"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// Encode as PNG
	if err := png.Encode(outputFile, img); err != nil {
		return fmt.Errorf("error encoding PNG: %v", err)
	}

	return nil
}
