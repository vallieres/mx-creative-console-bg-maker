package processor

import (
	"errors"
	"image"
	"image/color"
	"io"
)

const redColor = 255

// CreateTestImage creates a simple test image for use in tests.
func CreateTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c := color.RGBA{R: redColor, G: 0, B: 0, A: redColor} // Red

	for y := range height {
		for x := range width {
			img.Set(x, y, c)
		}
	}

	return img
}

// TestMockFileSystem implements FileSystem for testing across packages.
type TestMockFileSystem struct {
	OpenFunc   func(name string) (io.ReadCloser, error)
	CreateFunc func(name string) (io.WriteCloser, error)
	files      map[string][]byte
	written    map[string][]byte
}

// NewTestMockFileSystem creates a mock filesystem for testing.
func NewTestMockFileSystem() *TestMockFileSystem {
	return &TestMockFileSystem{
		files:   make(map[string][]byte),
		written: make(map[string][]byte),
	}
}

// AddFile adds a file to the mock filesystem.
func (m *TestMockFileSystem) AddFile(name string, content []byte) {
	m.files[name] = content
}

// GetWrittenFile returns the content written to a file.
func (m *TestMockFileSystem) GetWrittenFile(name string) ([]byte, bool) {
	content, exists := m.written[name]
	return content, exists
}

// Open implements FileSystem.
func (m *TestMockFileSystem) Open(name string) (io.ReadCloser, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(name)
	}

	content, exists := m.files[name]
	if !exists {
		return nil, errors.New("file not found")
	}
	return &testMockReadCloser{content: content}, nil
}

// Create implements FileSystem.
func (m *TestMockFileSystem) Create(name string) (io.WriteCloser, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(name)
	}

	return &testMockWriteCloser{
		fs:   m,
		name: name,
	}, nil
}

type testMockReadCloser struct {
	content []byte
	pos     int
}

func (m *testMockReadCloser) Read(p []byte) (int, error) {
	if m.pos >= len(m.content) {
		return 0, io.EOF
	}
	n := copy(p, m.content[m.pos:])
	m.pos += n
	return n, nil
}

func (m *testMockReadCloser) Close() error {
	return nil
}

type testMockWriteCloser struct {
	fs     *TestMockFileSystem
	name   string
	buffer []byte
}

func (m *testMockWriteCloser) Write(p []byte) (int, error) {
	m.buffer = append(m.buffer, p...)
	return len(p), nil
}

func (m *testMockWriteCloser) Close() error {
	m.fs.written[m.name] = m.buffer
	return nil
}

// TestMockImageDecoder implements ImageDecoder for testing across packages.
type TestMockImageDecoder struct {
	DecodeFunc func(r io.Reader) (image.Image, string, error)
	img        image.Image
	format     string
	err        error
}

// NewTestMockImageDecoder creates a new mock image decoder.
func NewTestMockImageDecoder(img image.Image, format string, err error) *TestMockImageDecoder {
	return &TestMockImageDecoder{img: img, format: format, err: err}
}

// Decode implements ImageDecoder.
func (m *TestMockImageDecoder) Decode(r io.Reader) (image.Image, string, error) {
	if m.DecodeFunc != nil {
		return m.DecodeFunc(r)
	}
	return m.img, m.format, m.err
}

// TestMockImageEncoder implements ImageEncoder for testing across packages.
type TestMockImageEncoder struct {
	EncodeFunc func(w io.Writer, img image.Image) error
	err        error
}

// NewTestMockImageEncoder creates a new mock image encoder.
func NewTestMockImageEncoder(err error) *TestMockImageEncoder {
	return &TestMockImageEncoder{err: err}
}

// Encode implements ImageEncoder.
func (m *TestMockImageEncoder) Encode(w io.Writer, img image.Image) error {
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

// TestMockImageResizer implements ImageResizer for testing across packages.
type TestMockImageResizer struct {
	ResizeFunc func(width, height uint, img image.Image) image.Image
}

// NewTestMockImageResizer creates a new mock image resizer.
func NewTestMockImageResizer() *TestMockImageResizer {
	return &TestMockImageResizer{}
}

// Resize implements ImageResizer.
func (m *TestMockImageResizer) Resize(width, height uint, img image.Image) image.Image {
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
	return CreateTestImage(int(width), int(height)) // #nosec G115
}

// CreateColoredTestImage creates a test image with a specific color.
func CreateColoredTestImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := range height {
		for x := range width {
			img.Set(x, y, c)
		}
	}

	return img
}
