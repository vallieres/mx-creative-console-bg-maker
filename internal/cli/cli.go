package cli

import (
	"fmt"
	"os"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

// App represents the CLI application.
type App struct {
	processor *processor.Service
}

// NewApp creates a new CLI application.
func NewApp() *App {
	return &App{
		processor: processor.NewService(),
	}
}

// NewAppWithProcessor creates a new CLI application with a custom processor.
func NewAppWithProcessor(proc *processor.Service) *App {
	return &App{
		processor: proc,
	}
}

const minRequiredArgs = 2

// Run executes the CLI application.
func (a *App) Run(args []string) error {
	progName := "ccbm"

	if len(args) >= 2 && (args[1] == "--help" || args[1] == "-h") {
		_, _ = fmt.Fprintf(os.Stdout, "Usage: ccbm <image_path>\n")
		return nil
	}

	if len(args) < minRequiredArgs {
		return fmt.Errorf("usage: %s <image_path>", progName)
	}

	imagePath := args[1]
	return a.processor.ProcessImage(imagePath)
}

// Main is the main entry point that can be tested.
func Main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
