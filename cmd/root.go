package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "manager",
	Long: "Job Manager for workers",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Run manager -h for instructions")
	},
}

// Execute commands.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
