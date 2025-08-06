package processor_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestDefaultConfig_Values(t *testing.T) {
	// Execute
	config := processor.DefaultConfig()

	// Assert
	assert.Equal(t, 378, config.TargetSize)
	assert.Equal(t, 3, config.GridSize)
	assert.Equal(t, 116, config.TileSize)
	assert.Equal(t, 15, config.Spacing)
}

func TestResizeImage_LandscapeImage(t *testing.T) {
	// Setup
	testImg := processor.CreateTestImage(800, 600)
	targetSize := 200
	var capturedWidth, capturedHeight uint

	resizer := processor.NewTestMockImageResizer()
	resizer.ResizeFunc = func(width, height uint, _ image.Image) image.Image {
		capturedWidth = width
		capturedHeight = height
		return processor.CreateTestImage(int(width), int(height))
	}

	// Execute
	result := processor.ResizeImage(testImg, targetSize, resizer)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, uint(0), capturedWidth, "Width should be 0 for landscape resize")
	assert.Equal(t, uint(targetSize), capturedHeight, "Height should be target size")
}

func TestResizeImage_PortraitImage(t *testing.T) {
	// Setup
	testImg := processor.CreateTestImage(600, 800)
	targetSize := 200
	var capturedWidth, capturedHeight uint

	resizer := processor.NewTestMockImageResizer()
	resizer.ResizeFunc = func(width, height uint, _ image.Image) image.Image {
		capturedWidth = width
		capturedHeight = height
		return processor.CreateTestImage(int(width), int(height))
	}

	// Execute
	result := processor.ResizeImage(testImg, targetSize, resizer)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, uint(targetSize), capturedWidth, "Width should be target size")
	assert.Equal(t, uint(0), capturedHeight, "Height should be 0 for portrait resize")
}

func TestResizeImage_SquareImage(t *testing.T) {
	// Setup
	testImg := processor.CreateTestImage(600, 600)
	targetSize := 200
	var capturedWidth, capturedHeight uint

	resizer := processor.NewTestMockImageResizer()
	resizer.ResizeFunc = func(width, height uint, _ image.Image) image.Image {
		capturedWidth = width
		capturedHeight = height
		return processor.CreateTestImage(int(width), int(height))
	}

	// Execute
	result := processor.ResizeImage(testImg, targetSize, resizer)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, uint(targetSize), capturedWidth, "Width should be target size")
	assert.Equal(t, uint(0), capturedHeight, "Height should be 0 for square image")
}

func TestCropToSquare_WideImage(t *testing.T) {
	// Setup
	testImg := processor.CreateTestImage(500, 300)
	targetSize := 200

	// Execute
	result := processor.CropToSquare(testImg, targetSize)

	// Assert
	assert.NotNil(t, result)
	bounds := result.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, targetSize, width)
	assert.Equal(t, targetSize, height)
}

func TestCropToSquare_TallImage(t *testing.T) {
	// Setup
	testImg := processor.CreateTestImage(300, 500)
	targetSize := 200

	// Execute
	result := processor.CropToSquare(testImg, targetSize)

	// Assert
	assert.NotNil(t, result)
	bounds := result.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, targetSize, width)
	assert.Equal(t, targetSize, height)
}

func TestCropToSquare_ExactSize(t *testing.T) {
	// Setup
	targetSize := 200
	testImg := processor.CreateTestImage(targetSize, targetSize)

	// Execute
	result := processor.CropToSquare(testImg, targetSize)

	// Assert
	assert.NotNil(t, result)
	bounds := result.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, targetSize, width)
	assert.Equal(t, targetSize, height)
}

func TestSplitIntoTiles_DefaultConfig(t *testing.T) {
	// Setup
	config := processor.DefaultConfig()
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	// Execute
	result := processor.SplitIntoTiles(testImg, config)

	// Assert
	expectedTileCount := config.GridSize * config.GridSize
	assert.Len(t, result.Tiles, expectedTileCount)
	assert.Len(t, result.TileCoords, expectedTileCount)

	// Verify each tile has correct dimensions
	for i, tile := range result.Tiles {
		bounds := tile.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y
		assert.Equal(t, config.TileSize, width, "Tile %d width incorrect", i)
		assert.Equal(t, config.TileSize, height, "Tile %d height incorrect", i)
	}

	// Verify tile coordinates are sequential and correct
	expectedCoords := []processor.TileCoordinate{
		{Row: 0, Col: 0, Number: 1},
		{Row: 0, Col: 1, Number: 2},
		{Row: 0, Col: 2, Number: 3},
		{Row: 1, Col: 0, Number: 4},
		{Row: 1, Col: 1, Number: 5},
		{Row: 1, Col: 2, Number: 6},
		{Row: 2, Col: 0, Number: 7},
		{Row: 2, Col: 1, Number: 8},
		{Row: 2, Col: 2, Number: 9},
	}

	assert.Equal(t, expectedCoords, result.TileCoords)
}

func TestSplitIntoTiles_CustomConfig(t *testing.T) {
	// Setup
	config := processor.Config{
		TargetSize: 200,
		GridSize:   2,
		TileSize:   90,
		Spacing:    10,
	}
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	// Execute
	result := processor.SplitIntoTiles(testImg, config)

	// Assert
	expectedTileCount := 4 // 2x2 grid
	assert.Len(t, result.Tiles, expectedTileCount)
	assert.Len(t, result.TileCoords, expectedTileCount)

	// Verify tile coordinates for 2x2 grid
	expectedCoords := []processor.TileCoordinate{
		{Row: 0, Col: 0, Number: 1},
		{Row: 0, Col: 1, Number: 2},
		{Row: 1, Col: 0, Number: 3},
		{Row: 1, Col: 1, Number: 4},
	}

	assert.Equal(t, expectedCoords, result.TileCoords)
}

func TestSplitIntoTiles_SingleTile(t *testing.T) {
	// Setup
	config := processor.Config{
		TargetSize: 100,
		GridSize:   1,
		TileSize:   90,
		Spacing:    5,
	}
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	// Execute
	result := processor.SplitIntoTiles(testImg, config)

	// Assert
	assert.Len(t, result.Tiles, 1)
	assert.Len(t, result.TileCoords, 1)

	expectedCoord := processor.TileCoordinate{Row: 0, Col: 0, Number: 1}
	assert.Equal(t, expectedCoord, result.TileCoords[0])
}

func TestSplitIntoTiles_LargeGrid(t *testing.T) {
	// Setup
	config := processor.Config{
		TargetSize: 500,
		GridSize:   5,
		TileSize:   90,
		Spacing:    10,
	}
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	// Execute
	result := processor.SplitIntoTiles(testImg, config)

	// Assert
	expectedTileCount := 25 // 5x5 grid
	assert.Len(t, result.Tiles, expectedTileCount)
	assert.Len(t, result.TileCoords, expectedTileCount)

	// Verify first and last coordinates
	firstCoord := processor.TileCoordinate{Row: 0, Col: 0, Number: 1}
	lastCoord := processor.TileCoordinate{Row: 4, Col: 4, Number: 25}
	assert.Equal(t, firstCoord, result.TileCoords[0])
	assert.Equal(t, lastCoord, result.TileCoords[24])
}

func TestTileCoordinate_NumberSequence(t *testing.T) {
	// Setup
	config := processor.Config{
		TargetSize: 300,
		GridSize:   3,
		TileSize:   90,
		Spacing:    10,
	}
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	// Execute
	result := processor.SplitIntoTiles(testImg, config)

	// Assert - verify numbers are sequential starting from 1
	for i, coord := range result.TileCoords {
		expectedNumber := i + 1
		assert.Equal(t, expectedNumber, coord.Number, "Tile number should be sequential")
	}
}

func TestProcessingResult_Structure(t *testing.T) {
	// Setup
	config := processor.DefaultConfig()
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	// Execute
	result := processor.SplitIntoTiles(testImg, config)

	// Assert structure integrity
	assert.IsType(t, processor.ProcessingResult{}, result)
	assert.IsType(t, []image.Image{}, result.Tiles)
	assert.IsType(t, []processor.TileCoordinate{}, result.TileCoords)

	// Verify tiles and coordinates have same length
	assert.Len(t, result.TileCoords, len(result.Tiles),
		"Tiles and coordinates should have same length")
}
