package cli_test

import (
	"errors"
	"io"
	"testing"

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

	if app == nil {
		t.Fatal("Expected non-nil app")
	}
}

func TestNewAppWithProcessor(t *testing.T) {
	customProcessor := processor.NewService()
	app := cli.NewAppWithProcessor(customProcessor)

	if app == nil {
		t.Fatal("Expected non-nil app")
	}
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
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
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

			if err == nil {
				t.Fatal("Expected error for missing arguments, got nil")
			}

			expectedMsg := "usage:"
			if !contains(err.Error(), expectedMsg) {
				t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
			}
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

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// Helper function to check if string contains substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
