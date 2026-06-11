package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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

// NewRootCmd builds the root cobra command with all subcommands attached.
func (a *App) NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ccbm",
		Short: "Command-line tool that splits images into a 3x3 grid for the Logitech MX Creative Console",
		Long: `
 ▄▄·  ▄▄· ▄▄▄▄· • ▌ ▄ ·.
▐█ ▌▪▐█ ▌▪▐█ ▀█▪·██ ▐███▪
██ ▄▄██ ▄▄▐█▀▀█▄▐█ ▌▐▌▐█·
▐███▌▐███▌██▄▪▐███ ██▌▐█▌
·▀▀▀ ·▀▀▀ ·▀▀▀▀ ▀▀  █▪▀▀▀
Creative Console BG Maker

The tool will create 9 PNG files in the same directory as the input image,
that represent the 9 keys on the MX Creative Console.
`,
		SilenceUsage: true,
	}

	var inputFile string

	splitCmd := &cobra.Command{
		Use:   "split",
		Short: "Breaks down the image in 9 squares (3x3)",
		Long: `Resizes and crops the input image to a square, then splits it into
9 tiles matching the 3x3 key layout of the MX Creative Console.
The tiles are saved as PNG files next to the input image.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return a.processor.ProcessImage(inputFile)
		},
	}

	splitCmd.Flags().StringVarP(&inputFile, "file", "f", "", "The input file")
	_ = splitCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(splitCmd)

	return rootCmd
}

// Run executes the CLI application. args is expected to be [os.Args]
// (the program name in args[0] is skipped).
func (a *App) Run(args []string) error {
	rootCmd := a.NewRootCmd()
	if len(args) > 1 {
		rootCmd.SetArgs(args[1:])
	} else {
		rootCmd.SetArgs([]string{})
	}
	return rootCmd.Execute()
}

// Main is the main entry point that can be tested.
func Main() {
	app := NewApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
