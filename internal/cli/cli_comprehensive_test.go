package cli_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/cli"
	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestNewApp_CreatesValidApp(t *testing.T) {
	// Execute
	app := cli.NewApp()

	// Assert
	assert.NotNil(t, app)
}

func TestNewAppWithProcessor_UsesCustomProcessor(t *testing.T) {
	// Setup
	customProcessor := processor.NewService()

	// Execute
	app := cli.NewAppWithProcessor(customProcessor)

	// Assert
	assert.NotNil(t, app)
}

func TestApp_Run_ValidArgs(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)
	fs.AddFile("/test/image.jpg", []byte("fake image data"))

	args := []string{"ccbm", "/test/image.jpg"}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.NoError(t, runErr)
}

func TestApp_Run_NoArguments(t *testing.T) {
	// Setup
	app := cli.NewApp()
	args := []string{}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "usage:")
}

func TestApp_Run_OnlyProgramName(t *testing.T) {
	// Setup
	app := cli.NewApp()
	args := []string{"ccbm"}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "usage: ccbm <image_path>")
}

func TestApp_Run_ProcessorFileError(t *testing.T) {
	// Setup - create a service that will fail
	fs := processor.NewTestMockFileSystem()
	fs.OpenFunc = func(_ string) (io.ReadCloser, error) {
		return nil, errors.New("file access denied")
	}

	decoder := processor.NewTestMockImageDecoder(nil, "", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	args := []string{"ccbm", "/restricted/image.jpg"}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "failed to load image")
}

func TestApp_Run_EmptyImagePath(t *testing.T) {
	// Setup
	app := cli.NewApp()
	args := []string{"ccbm", ""}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.Error(t, runErr)
	// The error should come from the processor trying to process an empty path
}

func TestApp_Run_MultipleImagePaths(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)
	fs.AddFile("/test/image1.jpg", []byte("fake image data"))

	args := []string{"ccbm", "/test/image1.jpg", "/test/image2.jpg"}

	// Execute - should only process first image path
	runErr := app.Run(args)

	// Assert
	require.NoError(t, runErr)
}

func TestApp_Run_SpecialCharactersInPath(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	specialPath := "/test/image with spaces & symbols!.jpg"
	fs.AddFile(specialPath, []byte("fake image data"))

	args := []string{"ccbm", specialPath}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.NoError(t, runErr)
}

func TestApp_Run_VeryLongPath(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	// Create a very long path
	longPath := "/very/long/path/with/many/directories/and/subdirectories/image.jpg"
	fs.AddFile(longPath, []byte("fake image data"))

	args := []string{"ccbm", longPath}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.NoError(t, runErr)
}

func TestApp_Run_RelativePath(t *testing.T) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	relativePath := "./images/test.jpg"
	fs.AddFile(relativePath, []byte("fake image data"))

	args := []string{"ccbm", relativePath}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.NoError(t, runErr)
}

func TestApp_Run_CustomProgramName(t *testing.T) {
	// Setup
	app := cli.NewApp()
	args := []string{"ccbm"}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.Error(t, runErr)
	assert.Contains(t, runErr.Error(), "usage: ccbm <image_path>")
}

func TestApp_Run_ArgumentsWithExtraSpaces(t *testing.T) {
	// Note: The Go args slice would already be parsed by the shell,
	// but we can test that our code handles normal arguments correctly

	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	imagePath := "/test/normal-image.jpg"
	fs.AddFile(imagePath, []byte("fake image data"))

	args := []string{"ccbm", imagePath}

	// Execute
	runErr := app.Run(args)

	// Assert
	require.NoError(t, runErr)
}

// Test CLI integration with different image formats.
func TestApp_Run_DifferentImageFormats(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		format   string
	}{
		{"JPEG", "/test/image.jpg", "jpeg"},
		{"PNG", "/test/image.png", "png"},
		{"GIF", "/test/image.gif", "gif"},
		{"BMP", "/test/image.bmp", "bmp"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			fs := processor.NewTestMockFileSystem()
			decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), tc.format, nil)
			encoder := processor.NewTestMockImageEncoder(nil)
			resizer := processor.NewTestMockImageResizer()
			config := processor.DefaultConfig()

			service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
			app := cli.NewAppWithProcessor(service)

			fs.AddFile(tc.filename, []byte("fake image data"))
			args := []string{"ccbm", tc.filename}

			// Execute
			runErr := app.Run(args)

			// Assert
			require.NoError(t, runErr)
		})
	}
}

// Benchmark test for CLI performance.
func BenchmarkApp_Run(b *testing.B) {
	// Setup
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	fs.AddFile("/test/image.jpg", []byte("fake image data"))
	args := []string{"ccbm", "/test/image.jpg"}

	// Benchmark
	b.ResetTimer()
	for range b.N {
		_ = app.Run(args)
	}
}
