package processor

import (
	"image"
	"io"
)

// FileSystem abstracts file system operations for testing.
type FileSystem interface {
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
}

// ImageDecoder abstracts image decoding operations.
type ImageDecoder interface {
	Decode(r io.Reader) (image.Image, string, error)
}

// ImageEncoder abstracts image encoding operations.
type ImageEncoder interface {
	Encode(w io.Writer, img image.Image) error
}

// ImageResizer abstracts image resizing operations.
type ImageResizer interface {
	Resize(width, height uint, img image.Image) image.Image
}
