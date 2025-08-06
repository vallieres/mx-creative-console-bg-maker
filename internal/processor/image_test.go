package processor_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestDefaultConfig(t *testing.T) {
	config := processor.DefaultConfig()

	expectedConfig := processor.Config{
		TargetSize: 378,
		GridSize:   3,
		TileSize:   116,
		Spacing:    15,
	}

	assert.Equal(t, expectedConfig, config)
}

func TestResizeImage(t *testing.T) {
	tests := []struct {
		name        string
		origWidth   int
		origHeight  int
		targetSize  int
		expectWidth bool // true if we expect width to be resized to 0
	}{
		{
			name:        "landscape image",
			origWidth:   800,
			origHeight:  600,
			targetSize:  378,
			expectWidth: true, // width > height, so resize height
		},
		{
			name:        "portrait image",
			origWidth:   600,
			origHeight:  800,
			targetSize:  378,
			expectWidth: false, // height > width, so resize width
		},
		{
			name:        "square image",
			origWidth:   600,
			origHeight:  600,
			targetSize:  378,
			expectWidth: false, // equal, so resize width
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test image
			testImg := processor.CreateTestImage(tt.origWidth, tt.origHeight)

			// Create mock resizer that tracks the parameters
			var capturedWidth, capturedHeight uint
			resizer := processor.NewTestMockImageResizer()
			resizer.ResizeFunc = func(width, height uint, _ image.Image) image.Image {
				capturedWidth = width
				capturedHeight = height
				return processor.CreateTestImage(int(width), int(height))
			}

			// Test the resize function
			result := processor.ResizeImage(testImg, tt.targetSize, resizer)

			// Verify the correct parameters were passed to the resizer
			if tt.expectWidth {
				assert.Equal(t, uint(0), capturedWidth, "Width should be 0 for landscape resize")
				assert.Equal(t, uint(tt.targetSize), capturedHeight, "Height should be target size")
			} else {
				assert.Equal(t, uint(tt.targetSize), capturedWidth, "Width should be target size")
				assert.Equal(t, uint(0), capturedHeight, "Height should be 0 for portrait/square resize")
			}

			// Verify we got an image back
			assert.NotNil(t, result)
		})
	}
}

func TestCropToSquare(t *testing.T) {
	tests := []struct {
		name       string
		imgWidth   int
		imgHeight  int
		targetSize int
	}{
		{
			name:       "crop wide image",
			imgWidth:   500,
			imgHeight:  300,
			targetSize: 200,
		},
		{
			name:       "crop tall image",
			imgWidth:   300,
			imgHeight:  500,
			targetSize: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testImg := processor.CreateTestImage(tt.imgWidth, tt.imgHeight)

			result := processor.CropToSquare(testImg, tt.targetSize)

			assert.NotNil(t, result)

			bounds := result.Bounds()
			width := bounds.Max.X - bounds.Min.X
			height := bounds.Max.Y - bounds.Min.Y

			assert.Equal(t, tt.targetSize, width)
			assert.Equal(t, tt.targetSize, height)
		})
	}
}

func TestSplitIntoTiles(t *testing.T) {
	config := processor.Config{
		TargetSize: 378,
		GridSize:   3,
		TileSize:   116,
		Spacing:    15,
	}

	// Create a test square image
	testImg := processor.CreateTestImage(config.TargetSize, config.TargetSize)

	result := processor.SplitIntoTiles(testImg, config)

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

	// Verify tile coordinates are correct
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
