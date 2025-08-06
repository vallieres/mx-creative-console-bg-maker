package processor

import (
	"fmt"
	"image"
	"path/filepath"
	"strings"
)

// Service handles image processing with configurable dependencies.
type Service struct {
	fileSystem FileSystem
	decoder    ImageDecoder
	encoder    ImageEncoder
	resizer    ImageResizer
	config     Config
}

// NewService creates a new processor service with default dependencies.
func NewService() *Service {
	return &Service{
		fileSystem: &OSFileSystem{},
		decoder:    &StandardImageDecoder{},
		encoder:    &PNGEncoder{},
		resizer:    &LanczosResizer{},
		config:     DefaultConfig(),
	}
}

// NewServiceWithDeps creates a new processor service with custom dependencies.
func NewServiceWithDeps(fs FileSystem, decoder ImageDecoder, encoder ImageEncoder, resizer ImageResizer, config Config) *Service {
	return &Service{
		fileSystem: fs,
		decoder:    decoder,
		encoder:    encoder,
		resizer:    resizer,
		config:     config,
	}
}

// ProcessImage processes an image file and splits it into tiles.
func (s *Service) ProcessImage(imagePath string) error {
	img, loadErr := s.LoadImage(imagePath)
	if loadErr != nil {
		return fmt.Errorf("failed to load image: %w", loadErr)
	}

	processedImg := s.ProcessImageData(img)

	if saveErr := s.SaveTiles(processedImg, imagePath); saveErr != nil {
		return fmt.Errorf("failed to save tiles: %w", saveErr)
	}

	return nil
}

// LoadImage loads and decodes an image from file.
func (s *Service) LoadImage(imagePath string) (*ProcessedImage, error) {
	file, openErr := s.fileSystem.Open(imagePath)
	if openErr != nil {
		return nil, fmt.Errorf("error opening image: %w", openErr)
	}
	defer file.Close()

	img, _, decodeErr := s.decoder.Decode(file)
	if decodeErr != nil {
		return nil, fmt.Errorf("error decoding image: %w", decodeErr)
	}

	return &ProcessedImage{Original: img}, nil
}

// ProcessedImage holds an image and its processed versions.
type ProcessedImage struct {
	Original image.Image
	Resized  image.Image
	Squared  image.Image
	Result   ProcessingResult
}

// ProcessImageData handles the core image processing logic.
func (s *Service) ProcessImageData(procImg *ProcessedImage) *ProcessedImage {
	procImg.Resized = ResizeImage(procImg.Original, s.config.TargetSize, s.resizer)
	procImg.Squared = CropToSquare(procImg.Resized, s.config.TargetSize)
	procImg.Result = SplitIntoTiles(procImg.Squared, s.config)
	return procImg
}

// SaveTiles saves all tiles to disk.
func (s *Service) SaveTiles(procImg *ProcessedImage, originalPath string) error {
	baseDir := filepath.Dir(originalPath)
	fileName := strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(originalPath))

	for i, tile := range procImg.Result.Tiles {
		coord := procImg.Result.TileCoords[i]
		outputPath := filepath.Join(baseDir, fmt.Sprintf("%s_%d.png", fileName, coord.Number))

		if saveErr := s.SaveTile(tile, outputPath); saveErr != nil {
			return fmt.Errorf("error saving tile %d: %w", coord.Number, saveErr)
		}
	}

	return nil
}

// SaveTile saves a single tile to disk.
func (s *Service) SaveTile(tile image.Image, outputPath string) error {
	outputFile, createErr := s.fileSystem.Create(outputPath)
	if createErr != nil {
		return fmt.Errorf("error creating output file: %w", createErr)
	}
	defer outputFile.Close()

	if encodeErr := s.encoder.Encode(outputFile, tile); encodeErr != nil {
		return fmt.Errorf("error encoding tile: %w", encodeErr)
	}

	return nil
}
