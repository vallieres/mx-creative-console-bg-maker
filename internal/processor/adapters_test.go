package processor_test

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestStandardImageDecoder_Decode_PNG(t *testing.T) {
	// Setup - create a real PNG image in memory
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(5, 5, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	var buf bytes.Buffer
	encodeErr := png.Encode(&buf, img)
	require.NoError(t, encodeErr)

	decoder := &processor.StandardImageDecoder{}

	// Execute
	decodedImg, format, decodeErr := decoder.Decode(&buf)

	// Assert
	require.NoError(t, decodeErr)
	assert.NotNil(t, decodedImg)
	assert.Equal(t, "png", format)

	// Verify the decoded image has the correct pixel
	decodedColor := decodedImg.At(5, 5)
	r, g, b, a := decodedColor.RGBA()
	assert.Equal(t, uint32(65535), r) // 255 * 257
	assert.Equal(t, uint32(0), g)
	assert.Equal(t, uint32(0), b)
	assert.Equal(t, uint32(65535), a) // 255 * 257
}

func TestStandardImageDecoder_Decode_JPEG(t *testing.T) {
	// Setup - create a real JPEG image in memory
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := range 10 {
		for x := range 10 {
			img.Set(x, y, color.RGBA{R: 128, G: 128, B: 128, A: 255})
		}
	}

	var buf bytes.Buffer
	encodeErr := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	require.NoError(t, encodeErr)

	decoder := &processor.StandardImageDecoder{}

	// Execute
	decodedImg, format, decodeErr := decoder.Decode(&buf)

	// Assert
	require.NoError(t, decodeErr)
	assert.NotNil(t, decodedImg)
	assert.Equal(t, "jpeg", format)

	bounds := decodedImg.Bounds()
	assert.Equal(t, 10, bounds.Max.X)
	assert.Equal(t, 10, bounds.Max.Y)
}

func TestStandardImageDecoder_Decode_InvalidData(t *testing.T) {
	// Setup
	decoder := &processor.StandardImageDecoder{}
	invalidData := strings.NewReader("not an image")

	// Execute
	decodedImg, format, decodeErr := decoder.Decode(invalidData)

	// Assert
	require.Error(t, decodeErr)
	assert.Nil(t, decodedImg)
	assert.Empty(t, format)
}

func TestPNGEncoder_Encode_Success(t *testing.T) {
	// Setup
	img := processor.CreateTestImage(50, 50)
	encoder := &processor.PNGEncoder{}
	var buf bytes.Buffer

	// Execute
	encodeErr := encoder.Encode(&buf, img)

	// Assert
	require.NoError(t, encodeErr)
	assert.Positive(t, buf.Len(), "PNG data should be written")

	// Verify we can decode the PNG back
	decoder := &processor.StandardImageDecoder{}
	decodedImg, format, decodeErr := decoder.Decode(&buf)
	require.NoError(t, decodeErr)
	assert.Equal(t, "png", format)
	assert.NotNil(t, decodedImg)

	// Verify dimensions
	bounds := decodedImg.Bounds()
	assert.Equal(t, 50, bounds.Max.X)
	assert.Equal(t, 50, bounds.Max.Y)
}

func TestPNGEncoder_Encode_NilImage(t *testing.T) {
	// Setup
	encoder := &processor.PNGEncoder{}
	var buf bytes.Buffer

	// Execute and Assert - this will panic with nil image, which is expected behavior
	// We test this to document the behavior, but we know it will panic
	assert.Panics(t, func() {
		encoder.Encode(&buf, nil)
	})
}

func TestLanczosResizer_Resize_Downscale(t *testing.T) {
	// Setup
	originalImg := processor.CreateTestImage(100, 100)
	resizer := &processor.LanczosResizer{}

	// Execute
	resizedImg := resizer.Resize(50, 50, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, 50, width)
	assert.Equal(t, 50, height)
}

func TestLanczosResizer_Resize_Upscale(t *testing.T) {
	// Setup
	originalImg := processor.CreateTestImage(50, 50)
	resizer := &processor.LanczosResizer{}

	// Execute
	resizedImg := resizer.Resize(100, 100, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, 100, width)
	assert.Equal(t, 100, height)
}

func TestLanczosResizer_Resize_MaintainAspectRatio(t *testing.T) {
	// Setup - create a rectangular image
	originalImg := processor.CreateTestImage(200, 100)
	resizer := &processor.LanczosResizer{}

	// Execute - resize with width=0 to maintain aspect ratio
	resizedImg := resizer.Resize(0, 50, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Width should be calculated to maintain 2:1 aspect ratio
	assert.Equal(t, 50, height)
	assert.Equal(t, 100, width) // Should be 2x the height
}

func TestLanczosResizer_Resize_MaintainAspectRatioHeight(t *testing.T) {
	// Setup - create a rectangular image
	originalImg := processor.CreateTestImage(100, 200)
	resizer := &processor.LanczosResizer{}

	// Execute - resize with height=0 to maintain aspect ratio
	resizedImg := resizer.Resize(50, 0, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Height should be calculated to maintain 1:2 aspect ratio
	assert.Equal(t, 50, width)
	assert.Equal(t, 100, height) // Should be 2x the width
}

func TestLanczosResizer_Resize_ZeroBoth(t *testing.T) {
	// Setup
	originalImg := processor.CreateTestImage(100, 100)
	resizer := &processor.LanczosResizer{}

	// Execute - both dimensions zero should return original size
	resizedImg := resizer.Resize(0, 0, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, 100, width)
	assert.Equal(t, 100, height)
}

func TestLanczosResizer_Resize_SinglePixel(t *testing.T) {
	// Setup
	originalImg := processor.CreateTestImage(100, 100)
	resizer := &processor.LanczosResizer{}

	// Execute
	resizedImg := resizer.Resize(1, 1, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, 1, width)
	assert.Equal(t, 1, height)
}

func TestLanczosResizer_Resize_LargeImage(t *testing.T) {
	// Setup
	originalImg := processor.CreateTestImage(1000, 1000)
	resizer := &processor.LanczosResizer{}

	// Execute
	resizedImg := resizer.Resize(100, 100, originalImg)

	// Assert
	assert.NotNil(t, resizedImg)
	bounds := resizedImg.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	assert.Equal(t, 100, width)
	assert.Equal(t, 100, height)
}

// Benchmark tests.
func BenchmarkPNGEncoder_Encode(b *testing.B) {
	img := processor.CreateTestImage(378, 378) // Default target size
	encoder := &processor.PNGEncoder{}

	b.ResetTimer()
	for range b.N {
		var buf bytes.Buffer
		_ = encoder.Encode(&buf, img)
	}
}

func BenchmarkLanczosResizer_Resize(b *testing.B) {
	originalImg := processor.CreateTestImage(800, 600)
	resizer := &processor.LanczosResizer{}

	b.ResetTimer()
	for range b.N {
		_ = resizer.Resize(400, 300, originalImg)
	}
}
