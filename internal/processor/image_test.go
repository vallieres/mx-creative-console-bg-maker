package processor_test

import (
	"image"
	"testing"

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

	if config != expectedConfig {
		t.Errorf("Expected config %+v, got %+v", expectedConfig, config)
	}
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
				if capturedWidth != 0 || capturedHeight != uint(tt.targetSize) {
					t.Errorf("Expected resize(0, %d), got resize(%d, %d)", tt.targetSize, capturedWidth, capturedHeight)
				}
			} else {
				if capturedWidth != uint(tt.targetSize) || capturedHeight != 0 {
					t.Errorf("Expected resize(%d, 0), got resize(%d, %d)", tt.targetSize, capturedWidth, capturedHeight)
				}
			}

			// Verify we got an image back
			if result == nil {
				t.Error("Expected non-nil result image")
			}
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

			if result == nil {
				t.Fatal("Expected non-nil result image")
			}

			bounds := result.Bounds()
			width := bounds.Max.X - bounds.Min.X
			height := bounds.Max.Y - bounds.Min.Y

			if width != tt.targetSize || height != tt.targetSize {
				t.Errorf("Expected %dx%d image, got %dx%d", tt.targetSize, tt.targetSize, width, height)
			}
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
	if len(result.Tiles) != expectedTileCount {
		t.Errorf("Expected %d tiles, got %d", expectedTileCount, len(result.Tiles))
	}

	if len(result.TileCoords) != expectedTileCount {
		t.Errorf("Expected %d tile coordinates, got %d", expectedTileCount, len(result.TileCoords))
	}

	// Verify each tile has correct dimensions
	for i, tile := range result.Tiles {
		bounds := tile.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y

		if width != config.TileSize || height != config.TileSize {
			t.Errorf("Tile %d: expected %dx%d, got %dx%d", i, config.TileSize, config.TileSize, width, height)
		}
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

	for i, coord := range result.TileCoords {
		if coord != expectedCoords[i] {
			t.Errorf("Tile %d: expected coordinate %+v, got %+v", i, expectedCoords[i], coord)
		}
	}
}
