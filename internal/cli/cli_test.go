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

// MockProcessor implements a mock processor for testing.
type MockProcessor struct {
	ProcessImageFunc func(imagePath string) error
	CallCount        int
	LastImagePath    string
}

func (m *MockProcessor) ProcessImage(imagePath string) error {
	m.CallCount++
	m.LastImagePath = imagePath
	if m.ProcessImageFunc != nil {
		return m.ProcessImageFunc(imagePath)
	}
	return nil
}

func TestNewApp(t *testing.T) {
	app := cli.NewApp()

	assert.NotNil(t, app)
}

func TestNewAppWithProcessor(t *testing.T) {
	customProcessor := processor.NewService()
	app := cli.NewAppWithProcessor(customProcessor)

	assert.NotNil(t, app)
}

func TestApp_Run_Success(t *testing.T) {
	// Create service with mock dependencies for testing
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(processor.CreateTestImage(100, 100), "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	args := []string{"ccbm", "/test/image.jpg"}

	// Add test file
	fs.AddFile("/test/image.jpg", []byte("fake image"))

	err := app.Run(args)
	assert.NoError(t, err)
}

func TestApp_Run_MissingArguments(t *testing.T) {
	app := cli.NewApp()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "only program name",
			args: []string{"ccbm"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.Run(tt.args)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "usage:")
		})
	}
}

func TestApp_Run_ProcessorError(t *testing.T) {
	// Create a service that will fail
	fs := processor.NewTestMockFileSystem()
	fs.OpenFunc = func(_ string) (io.ReadCloser, error) {
		return nil, errors.New("file not found")
	}

	decoder := processor.NewTestMockImageDecoder(nil, "", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)
	app := cli.NewAppWithProcessor(service)

	args := []string{"ccbm", "/nonexistent/image.jpg"}

	err := app.Run(args)

	require.Error(t, err)
}
