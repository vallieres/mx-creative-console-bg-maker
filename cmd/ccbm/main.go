package main

import (
	"log"
	"os"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

func main() {
	if len(os.Args) != 2 { //nolint:mnd
		log.Fatal("Usage: imgsplit <image_path>")
	}

	imagePath := os.Args[1]
	if err := processor.ProcessImage(imagePath); err != nil {
		log.Fatal(err)
	}
}
