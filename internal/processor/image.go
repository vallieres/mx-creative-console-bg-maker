package processor

import (
	"image"
	"image/draw"
)

const (
	defaultTargetSize = 378
	defaultGridSize   = 3
	defaultTileSize   = 116
	defaultSpacing    = 15
	centerDivisor     = 2
)

// Config holds the processing configuration.
type Config struct {
	TargetSize int
	GridSize   int
	TileSize   int
	Spacing    int
}

// DefaultConfig returns the default processing configuration.
func DefaultConfig() Config {
	return Config{
		TargetSize: defaultTargetSize,
		GridSize:   defaultGridSize,
		TileSize:   defaultTileSize,
		Spacing:    defaultSpacing,
	}
}

// ProcessingResult holds the result of image processing.
type ProcessingResult struct {
	Tiles      []image.Image
	TileCoords []TileCoordinate
}

// TileCoordinate represents the position of a tile in the grid.
type TileCoordinate struct {
	Row    int
	Col    int
	Number int
}

// ResizeImage resizes an image to fit within the target size.
func ResizeImage(img image.Image, targetSize int, resizer ImageResizer) image.Image {
	bounds := img.Bounds()
	origWidth := bounds.Max.X - bounds.Min.X
	origHeight := bounds.Max.Y - bounds.Min.Y

	targetSizeUint := uint(targetSize) // #nosec G115

	if origWidth > origHeight {
		return resizer.Resize(0, targetSizeUint, img)
	}
	return resizer.Resize(targetSizeUint, 0, img)
}

// CropToSquare crops an image to a square centered on the original.
func CropToSquare(img image.Image, targetSize int) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	squared := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	startX := (width - targetSize) / centerDivisor
	startY := (height - targetSize) / centerDivisor

	draw.Draw(squared, squared.Bounds(), img, image.Point{
		X: startX,
		Y: startY,
	}, draw.Src)

	return squared
}

// SplitIntoTiles splits a square image into a grid of tiles.
func SplitIntoTiles(img image.Image, config Config) ProcessingResult {
	var tiles []image.Image
	var coords []TileCoordinate

	count := 1
	for y := range config.GridSize {
		for x := range config.GridSize {
			srcX := x * (config.TileSize + config.Spacing)
			srcY := y * (config.TileSize + config.Spacing)

			tile := image.NewRGBA(image.Rect(0, 0, config.TileSize, config.TileSize))
			draw.Draw(tile, tile.Bounds(), img, image.Point{X: srcX, Y: srcY}, draw.Src)

			tiles = append(tiles, tile)
			coords = append(coords, TileCoordinate{
				Row:    y,
				Col:    x,
				Number: count,
			})
			count++
		}
	}

	return ProcessingResult{
		Tiles:      tiles,
		TileCoords: coords,
	}
}
