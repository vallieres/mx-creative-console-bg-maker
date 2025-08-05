package processor

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // Register JPEG format
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

const (
	targetSize = 378
	gridSize   = 3
	tileSize   = 116
	spacing    = 15
)

func ProcessImage(imagePath string) error {
	// Open and decode the image
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("error opening image: %w", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	origWidth := bounds.Max.X - bounds.Min.X
	origHeight := bounds.Max.Y - bounds.Min.Y

	// Original dimensions logged

	// Determine which dimension to resize based on
	var resized image.Image
	if origWidth > origHeight {
		resized = resize.Resize(0, targetSize, img, resize.Lanczos3)
	} else {
		resized = resize.Resize(targetSize, 0, img, resize.Lanczos3)
	}

	// Get new dimensions after resize
	newBounds := resized.Bounds()
	newWidth := newBounds.Max.X - newBounds.Min.X
	newHeight := newBounds.Max.Y - newBounds.Min.Y

	// Resize completed

	// Create a square image for cropping from the center
	squared := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	// Calculate the starting point for cropping (to center the image)
	startX := (newWidth - targetSize) / 2  //nolint:mnd
	startY := (newHeight - targetSize) / 2 //nolint:mnd

	// Draw the cropped portion
	draw.Draw(squared, squared.Bounds(), resized, image.Point{
		X: startX,
		Y: startY,
	}, draw.Src)

	// Split into 3x3 grid
	baseDir := filepath.Dir(imagePath)
	fileName := strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath))

	count := 1
	for y := range gridSize {
		for x := range gridSize {
			// Calculate positions with spacing
			srcX := x * (tileSize + spacing)
			srcY := y * (tileSize + spacing)

			// Create a new image for the tile
			tile := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
			draw.Draw(tile, tile.Bounds(), squared, image.Point{X: srcX, Y: srcY}, draw.Src)

			// Save the tile
			outputPath := filepath.Join(baseDir, fmt.Sprintf("%s_%d.png", fileName, count))
			outputFile, errCreate := os.Create(outputPath)
			if errCreate != nil {
				return fmt.Errorf("error creating output file %d: %w", count, errCreate)
			}

			if errEncode := png.Encode(outputFile, tile); errEncode != nil {
				_ = outputFile.Close() // Best effort cleanup
				return fmt.Errorf("error encoding tile %d: %w", count, errEncode)
			}

			if closeErr := outputFile.Close(); closeErr != nil {
				return fmt.Errorf("error closing output file %d: %w", count, closeErr)
			}
			// Tile created successfully
			count++
		}
	}

	return nil
}
