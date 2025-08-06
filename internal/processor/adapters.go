package processor

import (
	"image"
	"image/png"
	"io"

	"github.com/nfnt/resize"
)

// StandardImageDecoder implements ImageDecoder using Go's standard image package.
type StandardImageDecoder struct{}

// Decode decodes an image from a reader.
func (d *StandardImageDecoder) Decode(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

// PNGEncoder implements ImageEncoder for PNG format.
type PNGEncoder struct{}

// Encode encodes an image as PNG.
func (e *PNGEncoder) Encode(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

// LanczosResizer implements ImageResizer using Lanczos3 algorithm.
type LanczosResizer struct{}

// Resize resizes an image using Lanczos3 algorithm.
func (r *LanczosResizer) Resize(width, height uint, img image.Image) image.Image {
	return resize.Resize(width, height, img, resize.Lanczos3)
}
