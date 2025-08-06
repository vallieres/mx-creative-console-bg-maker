package processor_test

import (
	"errors"
	"fmt"
	"image"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestService_LoadImage_Success(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	testImg := processor.CreateTestImage(100, 100)
	decoder := processor.NewTestMockImageDecoder(testImg, "jpeg", nil)
	service := processor.NewServiceWithDeps(fs, decoder, nil, nil, processor.DefaultConfig())

	fs.AddFile("/test/image.jpg", []byte("fake jpeg data"))

	// Execute
	result, loadErr := service.LoadImage("/test/image.jpg")

	// Assert
	require.NoError(t, loadErr)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Original)
	assert.Equal(t, testImg, result.Original)
}

func TestService_LoadImage_FileNotFound(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	fs.OpenFunc = func(_ string) (io.ReadCloser, error) {
		return nil, errors.New("file not found")
	}
	decoder := processor.NewTestMockImageDecoder(nil, "", nil)
	service := processor.NewServiceWithDeps(fs, decoder, nil, nil, processor.DefaultConfig())

	// Execute
	result, loadErr := service.LoadImage("/nonexistent/image.jpg")

	// Assert
	require.Error(t, loadErr)
	assert.Nil(t, result)
	assert.Contains(t, loadErr.Error(), "error opening image")
}

func TestService_LoadImage_DecodeError(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(nil, "", errors.New("invalid image format"))
	service := processor.NewServiceWithDeps(fs, decoder, nil, nil, processor.DefaultConfig())

	fs.AddFile("/test/corrupt.jpg", []byte("corrupt data"))

	// Execute
	result, loadErr := service.LoadImage("/test/corrupt.jpg")

	// Assert
	require.Error(t, loadErr)
	assert.Nil(t, result)
	assert.Contains(t, loadErr.Error(), "error decoding image")
}

func TestService_ProcessImageData(t *testing.T) {
	// Setup
	resizer := processor.NewTestMockImageResizer()
	config := processor.Config{
		TargetSize: 200,
		GridSize:   2,
		TileSize:   90,
		Spacing:    10,
	}
	service := processor.NewServiceWithDeps(nil, nil, nil, resizer, config)
	testImg := processor.CreateTestImage(300, 200)
	procImg := &processor.ProcessedImage{Original: testImg}

	// Execute
	result := service.ProcessImageData(procImg)

	// Assert
	assert.NotNil(t, result)
	assert.NotNil(t, result.Resized)
	assert.NotNil(t, result.Squared)
	assert.Len(t, result.Result.Tiles, 4) // 2x2 grid
	assert.Len(t, result.Result.TileCoords, 4)
}

func TestService_SaveTile_Success(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	encoder := processor.NewTestMockImageEncoder(nil)
	service := processor.NewServiceWithDeps(fs, nil, encoder, nil, processor.DefaultConfig())
	testImg := processor.CreateTestImage(100, 100)

	// Execute
	saveErr := service.SaveTile(testImg, "/test/tile.png")

	// Assert
	require.NoError(t, saveErr)
	data, exists := fs.GetWrittenFile("/test/tile.png")
	assert.True(t, exists)
	assert.Equal(t, "fake png data", string(data))
}

func TestService_SaveTile_CreateFileError(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	fs.CreateFunc = func(_ string) (io.WriteCloser, error) {
		return nil, errors.New("permission denied")
	}
	encoder := processor.NewTestMockImageEncoder(nil)
	service := processor.NewServiceWithDeps(fs, nil, encoder, nil, processor.DefaultConfig())
	testImg := processor.CreateTestImage(100, 100)

	// Execute
	saveErr := service.SaveTile(testImg, "/test/tile.png")

	// Assert
	require.Error(t, saveErr)
	assert.Contains(t, saveErr.Error(), "error creating output file")
}

func TestService_SaveTile_EncodeError(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	encoder := processor.NewTestMockImageEncoder(errors.New("encoding failed"))
	service := processor.NewServiceWithDeps(fs, nil, encoder, nil, processor.DefaultConfig())
	testImg := processor.CreateTestImage(100, 100)

	// Execute
	saveErr := service.SaveTile(testImg, "/test/tile.png")

	// Assert
	require.Error(t, saveErr)
	assert.Contains(t, saveErr.Error(), "error encoding tile")
}

func TestService_SaveTiles_Success(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	encoder := processor.NewTestMockImageEncoder(nil)
	service := processor.NewServiceWithDeps(fs, nil, encoder, nil, processor.DefaultConfig())

	// Create processed image with 2x2 grid of tiles
	procImg := &processor.ProcessedImage{
		Result: processor.ProcessingResult{
			Tiles: []image.Image{
				processor.CreateTestImage(50, 50),
				processor.CreateTestImage(50, 50),
				processor.CreateTestImage(50, 50),
				processor.CreateTestImage(50, 50),
			},
			TileCoords: []processor.TileCoordinate{
				{Row: 0, Col: 0, Number: 1},
				{Row: 0, Col: 1, Number: 2},
				{Row: 1, Col: 0, Number: 3},
				{Row: 1, Col: 1, Number: 4},
			},
		},
	}

	// Execute
	saveErr := service.SaveTiles(procImg, "/test/original.jpg")

	// Assert
	require.NoError(t, saveErr)

	expectedFiles := []string{
		"/test/original_1.png",
		"/test/original_2.png",
		"/test/original_3.png",
		"/test/original_4.png",
	}

	for _, filename := range expectedFiles {
		data, exists := fs.GetWrittenFile(filename)
		assert.True(t, exists, "File should exist: %s", filename)
		assert.Equal(t, "fake png data", string(data))
	}
}

func TestService_SaveTiles_TileError(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	encoder := processor.NewTestMockImageEncoder(errors.New("encoding failed"))
	service := processor.NewServiceWithDeps(fs, nil, encoder, nil, processor.DefaultConfig())

	procImg := &processor.ProcessedImage{
		Result: processor.ProcessingResult{
			Tiles: []image.Image{processor.CreateTestImage(50, 50)},
			TileCoords: []processor.TileCoordinate{
				{Row: 0, Col: 0, Number: 1},
			},
		},
	}

	// Execute
	saveErr := service.SaveTiles(procImg, "/test/original.jpg")

	// Assert
	require.Error(t, saveErr)
	assert.Contains(t, saveErr.Error(), "error saving tile 1")
}

func TestService_ProcessImage_FullWorkflow(t *testing.T) {
	// Setup complete mock environment
	fs := processor.NewTestMockFileSystem()
	testImg := processor.CreateTestImage(400, 300)
	decoder := processor.NewTestMockImageDecoder(testImg, "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	fs.AddFile("/test/image.jpg", []byte("fake jpeg data"))

	// Execute
	processErr := service.ProcessImage("/test/image.jpg")

	// Assert
	require.NoError(t, processErr)

	// Verify all 9 tiles were created (3x3 grid)
	for i := 1; i <= 9; i++ {
		filename := fmt.Sprintf("/test/image_%d.png", i)
		data, exists := fs.GetWrittenFile(filename)
		assert.True(t, exists, "Tile file should exist: %s", filename)
		assert.Equal(t, "fake png data", string(data))
	}
}
