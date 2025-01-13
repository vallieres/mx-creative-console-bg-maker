package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/vallieres/mx-creative-console-bg-maker/internal/processor"
)

// splitCmd represents the split command.
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Breaks down the image in 9 squares (3x3)",
	Long: `

`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imagePath := args[0]
		fmt.Println("Source: ", imagePath)
		if err := processor.ProcessImage(imagePath); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// splitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// splitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
