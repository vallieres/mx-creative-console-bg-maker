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
	img, err := s.loadImage(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	processedImg := s.processImage(img)

	if saveErr := s.saveTiles(processedImg, imagePath); saveErr != nil {
		return fmt.Errorf("failed to save tiles: %w", saveErr)
	}

	return nil
}

// loadImage loads and decodes an image from file.
func (s *Service) loadImage(imagePath string) (*ProcessedImage, error) {
	file, err := s.fileSystem.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening image: %w", err)
	}
	defer file.Close()

	img, _, err := s.decoder.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
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

// processImage handles the core image processing logic.
func (s *Service) processImage(procImg *ProcessedImage) *ProcessedImage {
	procImg.Resized = ResizeImage(procImg.Original, s.config.TargetSize, s.resizer)
	procImg.Squared = CropToSquare(procImg.Resized, s.config.TargetSize)
	procImg.Result = SplitIntoTiles(procImg.Squared, s.config)
	return procImg
}

// saveTiles saves all tiles to disk.
func (s *Service) saveTiles(procImg *ProcessedImage, originalPath string) error {
	baseDir := filepath.Dir(originalPath)
	fileName := strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(originalPath))

	for i, tile := range procImg.Result.Tiles {
		coord := procImg.Result.TileCoords[i]
		outputPath := filepath.Join(baseDir, fmt.Sprintf("%s_%d.png", fileName, coord.Number))

		if err := s.saveTile(tile, outputPath); err != nil {
			return fmt.Errorf("error saving tile %d: %w", coord.Number, err)
		}
	}

	return nil
}

// saveTile saves a single tile to disk.
func (s *Service) saveTile(tile image.Image, outputPath string) error {
	outputFile, err := s.fileSystem.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	if encodeErr := s.encoder.Encode(outputFile, tile); encodeErr != nil {
		return fmt.Errorf("error encoding tile: %w", encodeErr)
	}

	return nil
}
