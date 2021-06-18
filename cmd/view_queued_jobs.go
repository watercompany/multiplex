package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/watercompany/multiplex/job"
)

var (
	showAll bool
)

func init() {
	rootCmd.AddCommand(viewJobsCmd)
	viewJobsCmd.Flags().StringVar(
		&databaseAddr,
		"db-addr",
		"",
		"address of database to save jobs",
	)
	viewJobsCmd.Flags().BoolVar(
		&showAll,
		"show-all",
		false,
		"show all queued jobs",
	)
}

var viewJobsCmd = &cobra.Command{
	Use:  "view-queued-jobs",
	Long: "See all queued jobs",
	Run: func(cmd *cobra.Command, args []string) {
		kv, err := job.ListAllJobs(databaseAddr)
		if err != nil {
			fmt.Printf("unable to add job: %v", err)
			return
		}

		if len(kv) == 0 {
			fmt.Printf("Currently no queued jobs.")
			return
		}

		fmt.Printf("All queued jobs:\n")

		if showAll {
			for key, val := range kv {
				fmt.Printf("%v:%v\n", key, val)
			}
		} else {
			fmt.Printf("%v:%v\n", "Queued jobs", kv["job-last-index"])
			fmt.Printf("%v:%v\n", "Queued jobs", kv["job-last-index"])
		}

	},
}
