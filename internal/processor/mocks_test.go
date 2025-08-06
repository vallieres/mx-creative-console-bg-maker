package processor_test

import (
	"errors"
	"image"
	"io"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

// MockFileSystem implements processor.FileSystem for testing.
type MockFileSystem struct {
	OpenFunc   func(name string) (io.ReadCloser, error)
	CreateFunc func(name string) (io.WriteCloser, error)
	files      map[string][]byte
	written    map[string][]byte
}

// NewMockFileSystem creates a new mock file system.
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files:   make(map[string][]byte),
		written: make(map[string][]byte),
	}
}

// AddFile adds a file to the mock filesystem.
func (m *MockFileSystem) AddFile(name string, content []byte) {
	m.files[name] = content
}

// GetWrittenFile returns the content written to a file.
func (m *MockFileSystem) GetWrittenFile(name string) ([]byte, bool) {
	content, exists := m.written[name]
	return content, exists
}

// Open implements processor.FileSystem.
func (m *MockFileSystem) Open(name string) (io.ReadCloser, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(name)
	}

	content, exists := m.files[name]
	if !exists {
		return nil, errors.New("file not found")
	}
	return &mockReadCloser{content: content}, nil
}

// Create implements processor.FileSystem.
func (m *MockFileSystem) Create(name string) (io.WriteCloser, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(name)
	}

	return &mockWriteCloser{
		fs:   m,
		name: name,
	}, nil
}

type mockReadCloser struct {
	content []byte
	pos     int
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	if m.pos >= len(m.content) {
		return 0, io.EOF
	}
	n := copy(p, m.content[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockReadCloser) Close() error {
	return nil
}

type mockWriteCloser struct {
	fs     *MockFileSystem
	name   string
	buffer []byte
}

func (m *mockWriteCloser) Write(p []byte) (int, error) {
	m.buffer = append(m.buffer, p...)
	return len(p), nil
}

func (m *mockWriteCloser) Close() error {
	m.fs.written[m.name] = m.buffer
	return nil
}

// MockImageDecoder implements processor.ImageDecoder for testing.
type MockImageDecoder struct {
	DecodeFunc func(r io.Reader) (image.Image, string, error)
	img        image.Image
	format     string
	err        error
}

// NewMockImageDecoder creates a new mock image decoder.
func NewMockImageDecoder(img image.Image, format string, err error) *MockImageDecoder {
	return &MockImageDecoder{img: img, format: format, err: err}
}

// Decode implements processor.ImageDecoder.
func (m *MockImageDecoder) Decode(r io.Reader) (image.Image, string, error) {
	if m.DecodeFunc != nil {
		return m.DecodeFunc(r)
	}
	return m.img, m.format, m.err
}

// MockImageEncoder implements processor.ImageEncoder for testing.
type MockImageEncoder struct {
	EncodeFunc func(w io.Writer, img image.Image) error
	err        error
}

// NewMockImageEncoder creates a new mock image encoder.
func NewMockImageEncoder(err error) *MockImageEncoder {
	return &MockImageEncoder{err: err}
}

// Encode implements processor.ImageEncoder.
func (m *MockImageEncoder) Encode(w io.Writer, img image.Image) error {
	if m.EncodeFunc != nil {
		return m.EncodeFunc(w, img)
	}
	if m.err != nil {
		return m.err
	}
	// Write some dummy data
	_, err := w.Write([]byte("fake png data"))
	return err
}

// MockImageResizer implements processor.ImageResizer for testing.
type MockImageResizer struct {
	ResizeFunc func(width, height uint, img image.Image) image.Image
}

// NewMockImageResizer creates a new mock image resizer.
func NewMockImageResizer() *MockImageResizer {
	return &MockImageResizer{}
}

// Resize implements processor.ImageResizer.
func (m *MockImageResizer) Resize(width, height uint, img image.Image) image.Image {
	if m.ResizeFunc != nil {
		return m.ResizeFunc(width, height, img)
	}
	// Return a simple colored rectangle for testing
	if width == 0 {
		width = height
	}
	if height == 0 {
		height = width
	}
	return processor.CreateTestImage(int(width), int(height))
}
