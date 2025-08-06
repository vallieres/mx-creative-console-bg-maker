package processor

import (
	"io"
	"os"
)

// OSFileSystem implements FileSystem using the actual OS file system.
type OSFileSystem struct{}

// Open opens a file for reading.
func (fs *OSFileSystem) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// Create creates a file for writing.
func (fs *OSFileSystem) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}
