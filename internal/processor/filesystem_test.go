package processor_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func TestOSFileSystem_Open_Success(t *testing.T) {
	// Setup - create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	writeErr := os.WriteFile(tempFile, []byte(testContent), 0o644)
	require.NoError(t, writeErr)

	fs := &processor.OSFileSystem{}

	// Execute
	file, openErr := fs.Open(tempFile)

	// Assert
	require.NoError(t, openErr)
	assert.NotNil(t, file)

	// Verify content
	content := make([]byte, len(testContent))
	n, readErr := file.Read(content)
	require.NoError(t, readErr)
	assert.Equal(t, len(testContent), n)
	assert.Equal(t, testContent, string(content))

	// Clean up
	closeErr := file.Close()
	require.NoError(t, closeErr)
}

func TestOSFileSystem_Open_FileNotFound(t *testing.T) {
	// Setup
	fs := &processor.OSFileSystem{}
	nonExistentFile := "/path/that/does/not/exist/file.txt"

	// Execute
	file, openErr := fs.Open(nonExistentFile)

	// Assert
	require.Error(t, openErr)
	assert.Nil(t, file)
	assert.Contains(t, openErr.Error(), "no such file or directory")
}

func TestOSFileSystem_Open_Directory(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	fs := &processor.OSFileSystem{}

	// Execute - try to open a directory as a file
	file, openErr := fs.Open(tempDir)

	// Assert
	// Opening a directory should succeed but reading from it might behave differently
	require.NoError(t, openErr)
	assert.NotNil(t, file)

	closeErr := file.Close()
	require.NoError(t, closeErr)
}

func TestOSFileSystem_Create_Success(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	newFile := filepath.Join(tempDir, "created.txt")
	fs := &processor.OSFileSystem{}
	testContent := "created content"

	// Execute
	file, createErr := fs.Create(newFile)

	// Assert
	require.NoError(t, createErr)
	assert.NotNil(t, file)

	// Write and verify
	n, writeErr := file.Write([]byte(testContent))
	require.NoError(t, writeErr)
	assert.Equal(t, len(testContent), n)

	closeErr := file.Close()
	require.NoError(t, closeErr)

	// Verify file was created
	content, readErr := os.ReadFile(newFile)
	require.NoError(t, readErr)
	assert.Equal(t, testContent, string(content))
}

func TestOSFileSystem_Create_InvalidPath(t *testing.T) {
	// Setup
	fs := &processor.OSFileSystem{}
	invalidPath := "/root/cannot/create/here.txt" // Assuming no write permission

	// Execute
	file, createErr := fs.Create(invalidPath)

	// Assert
	require.Error(t, createErr)
	assert.Nil(t, file)
}

func TestOSFileSystem_Create_OverwriteExisting(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "existing.txt")
	originalContent := "original content"
	newContent := "new content"

	// Create original file
	writeErr := os.WriteFile(existingFile, []byte(originalContent), 0o644)
	require.NoError(t, writeErr)

	fs := &processor.OSFileSystem{}

	// Execute - create should overwrite
	file, createErr := fs.Create(existingFile)

	// Assert
	require.NoError(t, createErr)
	assert.NotNil(t, file)

	// Write new content
	n, writeErr2 := file.Write([]byte(newContent))
	require.NoError(t, writeErr2)
	assert.Equal(t, len(newContent), n)

	closeErr := file.Close()
	require.NoError(t, closeErr)

	// Verify file was overwritten
	content, readErr := os.ReadFile(existingFile)
	require.NoError(t, readErr)
	assert.Equal(t, newContent, string(content))
}

func TestOSFileSystem_Create_NestedDirectory(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	createErr := os.MkdirAll(subDir, 0o755)
	require.NoError(t, createErr)

	nestedFile := filepath.Join(subDir, "nested.txt")
	fs := &processor.OSFileSystem{}
	testContent := "nested content"

	// Execute
	file, createFileErr := fs.Create(nestedFile)

	// Assert
	require.NoError(t, createFileErr)
	assert.NotNil(t, file)

	// Write content
	n, writeErr := file.Write([]byte(testContent))
	require.NoError(t, writeErr)
	assert.Equal(t, len(testContent), n)

	closeErr := file.Close()
	require.NoError(t, closeErr)

	// Verify file exists
	_, statErr := os.Stat(nestedFile)
	require.NoError(t, statErr)
}

func TestOSFileSystem_Create_EmptyFile(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.txt")
	fs := &processor.OSFileSystem{}

	// Execute
	file, createErr := fs.Create(emptyFile)

	// Assert
	require.NoError(t, createErr)
	assert.NotNil(t, file)

	// Close immediately without writing
	closeErr := file.Close()
	require.NoError(t, closeErr)

	// Verify empty file exists
	info, statErr := os.Stat(emptyFile)
	require.NoError(t, statErr)
	assert.Equal(t, int64(0), info.Size())
}

func TestOSFileSystem_InterfaceCompliance(t *testing.T) {
	// This test ensures OSFileSystem implements the FileSystem interface
	var fs processor.FileSystem = &processor.OSFileSystem{}
	assert.NotNil(t, fs)
}

func TestOSFileSystem_RoundTrip(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "roundtrip.txt")
	fs := &processor.OSFileSystem{}
	testContent := "round trip test content"

	// Execute - Create and write
	writeFile, createErr := fs.Create(testFile)
	require.NoError(t, createErr)

	n, writeErr := writeFile.Write([]byte(testContent))
	require.NoError(t, writeErr)
	assert.Equal(t, len(testContent), n)

	closeWriteErr := writeFile.Close()
	require.NoError(t, closeWriteErr)

	// Execute - Open and read
	readFile, openErr := fs.Open(testFile)
	require.NoError(t, openErr)

	content := make([]byte, len(testContent))
	readN, readErr := readFile.Read(content)
	require.NoError(t, readErr)
	assert.Equal(t, len(testContent), readN)

	closeReadErr := readFile.Close()
	require.NoError(t, closeReadErr)

	// Assert
	assert.Equal(t, testContent, string(content))
}

// Test with special file names.
func TestOSFileSystem_SpecialCharacters(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	specialFile := filepath.Join(tempDir, "file with spaces & symbols!.txt")
	fs := &processor.OSFileSystem{}
	testContent := "special characters test"

	// Execute
	file, createErr := fs.Create(specialFile)
	require.NoError(t, createErr)

	_, writeErr := file.Write([]byte(testContent))
	require.NoError(t, writeErr)

	closeErr := file.Close()
	require.NoError(t, closeErr)

	// Verify we can read it back
	readFile, openErr := fs.Open(specialFile)
	require.NoError(t, openErr)

	content := make([]byte, len(testContent))
	_, readErr := readFile.Read(content)
	require.NoError(t, readErr)

	closeReadErr := readFile.Close()
	require.NoError(t, closeReadErr)

	// Assert
	assert.Equal(t, testContent, string(content))
}

// Benchmark tests.
func BenchmarkOSFileSystem_Create(b *testing.B) {
	tempDir := b.TempDir()
	fs := &processor.OSFileSystem{}
	testContent := []byte("benchmark test content")

	b.ResetTimer()
	for i := range b.N {
		testFile := filepath.Join(tempDir, fmt.Sprintf("bench_%d.txt", i))
		file, _ := fs.Create(testFile)
		file.Write(testContent)
		file.Close()
	}
}

func BenchmarkOSFileSystem_Open(b *testing.B) {
	// Setup
	tempDir := b.TempDir()
	testFile := filepath.Join(tempDir, "bench_read.txt")
	testContent := []byte("benchmark read content")

	os.WriteFile(testFile, testContent, 0o644)

	fs := &processor.OSFileSystem{}

	b.ResetTimer()
	for range b.N {
		file, _ := fs.Open(testFile)
		content := make([]byte, len(testContent))
		file.Read(content)
		file.Close()
	}
}
