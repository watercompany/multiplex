package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/watercompany/multiplex/job"
)

func init() {
	rootCmd.AddCommand(viewJobsCmd)
}

var viewJobsCmd = &cobra.Command{
	Use:  "view-jobs",
	Long: "See all queued jobs",
	Run: func(cmd *cobra.Command, args []string) {
		kv, err := job.ListAllJobs()
		if err != nil {
			fmt.Printf("unable to add job: %v", err)
			return
		}

		if len(kv) == 0 {
			fmt.Printf("Currently no queued jobs.")
			return
		}

		fmt.Printf("All queued jobs:\n")
		for key, val := range kv {
			fmt.Printf("%v:%v\n", key, val)
		}
	},
}
