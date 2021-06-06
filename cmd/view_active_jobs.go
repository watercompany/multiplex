package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/watercompany/multiplex/job"
)

func init() {
	rootCmd.AddCommand(viewActiveJobsCmd)
	viewActiveJobsCmd.Flags().StringVar(
		&databaseAddr,
		"db-addr",
		"",
		"address of database to save jobs",
	)
}

var viewActiveJobsCmd = &cobra.Command{
	Use:  "view-active-jobs",
	Long: "See all active jobs",
	Run: func(cmd *cobra.Command, args []string) {
		actJobs, err := job.GetNumberOfActiveJobs(databaseAddr)
		if err != nil {
			fmt.Printf("unable to get current active jobs: %v", err)
			return
		}

		fmt.Printf("Current number of active jobs: %v\n", actJobs)
	},
}
