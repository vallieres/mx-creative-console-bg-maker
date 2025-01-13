package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "mx-creative-console-bg-maker",
	Short: "Command-line tool that splits images into a 3x3 grid for the Logitech MX Creative Console",
	Long: `
 ▄▄·  ▄▄· ▄▄▄▄· • ▌ ▄ ·.
▐█ ▌▪▐█ ▌▪▐█ ▀█▪·██ ▐███▪
██ ▄▄██ ▄▄▐█▀▀█▄▐█ ▌▐▌▐█·
▐███▌▐███▌██▄▪▐███ ██▌▐█▌
·▀▀▀ ·▀▀▀ ·▀▀▀▀ ▀▀  █▪▀▀▀
Creative Console BG Maker

The tool will create 9 PNG files in the same directory as the input image,
that represents the 9 keys on the MX Creative Console.".
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceUsage = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "More verbose output")
}
