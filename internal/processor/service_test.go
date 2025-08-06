package processor_test

import (
	"errors"
	"io"
	"testing"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestNewService(t *testing.T) {
	service := processor.NewService()

	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestNewServiceWithDeps(t *testing.T) {
	fs := processor.NewTestMockFileSystem()
	decoder := processor.NewTestMockImageDecoder(nil, "", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestService_ProcessImage_Success(t *testing.T) {
	// Setup mocks
	fs := processor.NewTestMockFileSystem()
	testImg := processor.CreateTestImage(400, 300)
	decoder := processor.NewTestMockImageDecoder(testImg, "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	// Add a fake image file
	fs.AddFile("/test/image.jpg", []byte("fake jpeg data"))

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)

	err := service.ProcessImage("/test/image.jpg")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify that 9 tile files were created
	expectedFiles := []string{
		"/test/image_1.png",
		"/test/image_2.png",
		"/test/image_3.png",
		"/test/image_4.png",
		"/test/image_5.png",
		"/test/image_6.png",
		"/test/image_7.png",
		"/test/image_8.png",
		"/test/image_9.png",
	}

	for _, filename := range expectedFiles {
		_, exists := fs.GetWrittenFile(filename)
		if !exists {
			t.Errorf("Expected file %s to be created", filename)
		}
	}
}

func TestService_ProcessImage_FileOpenError(t *testing.T) {
	fs := processor.NewTestMockFileSystem()
	fs.OpenFunc = func(_ string) (io.ReadCloser, error) {
		return nil, errors.New("file not found")
	}

	decoder := processor.NewTestMockImageDecoder(nil, "", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)

	err := service.ProcessImage("/nonexistent/image.jpg")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedMsg := "failed to load image"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

func TestService_ProcessImage_DecodeError(t *testing.T) {
	fs := processor.NewTestMockFileSystem()
	fs.AddFile("/test/image.jpg", []byte("fake jpeg data"))

	decoder := processor.NewTestMockImageDecoder(nil, "", errors.New("invalid image format"))
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)

	err := service.ProcessImage("/test/image.jpg")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedMsg := "failed to load image"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

func TestService_ProcessImage_EncodeError(t *testing.T) {
	fs := processor.NewTestMockFileSystem()
	testImg := processor.CreateTestImage(400, 300)
	decoder := processor.NewTestMockImageDecoder(testImg, "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(errors.New("encoding failed"))
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	fs.AddFile("/test/image.jpg", []byte("fake jpeg data"))

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)

	err := service.ProcessImage("/test/image.jpg")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedMsg := "failed to save tiles"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

func TestService_ProcessImage_CreateFileError(t *testing.T) {
	fs := processor.NewTestMockFileSystem()
	fs.CreateFunc = func(_ string) (io.WriteCloser, error) {
		return nil, errors.New("permission denied")
	}

	testImg := processor.CreateTestImage(400, 300)
	decoder := processor.NewTestMockImageDecoder(testImg, "jpeg", nil)
	encoder := processor.NewTestMockImageEncoder(nil)
	resizer := processor.NewTestMockImageResizer()
	config := processor.DefaultConfig()

	fs.AddFile("/test/image.jpg", []byte("fake jpeg data"))

	service := processor.NewServiceWithDeps(fs, decoder, encoder, resizer, config)

	err := service.ProcessImage("/test/image.jpg")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedMsg := "failed to save tiles"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error containing '%s', got: %v", expectedMsg, err)
	}
}

// Skip testing loadImage since it's an internal method and covered by ProcessImage tests

// Skip testing processImage since it's an internal method and covered by ProcessImage tests

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
