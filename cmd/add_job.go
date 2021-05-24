package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/watercompany/multiplex/job"
	"github.com/watercompany/multiplex/worker"
	"github.com/watercompany/multiplex/worker/client"
)

var (
	LogName  string
	TaskName string
)

func init() {
	rootCmd.AddCommand(addJobCmd)
	addJobCmd.Flags().StringVar(
		&LogName,
		"log-name",
		"test",
		"output log name",
	)

	addJobCmd.Flags().StringVar(
		&TaskName,
		"task",
		"pos",
		"name of the task the worker will do",
	)
}

var addJobCmd = &cobra.Command{
	Use:  "add-job",
	Long: "Adds new job",
	Run: func(cmd *cobra.Command, args []string) {
		var wCfg worker.WorkerCfg
		_, err := wCfg.GetWorkerCfg()
		if err != nil {
			panic(err)
		}

		var addArgs []string
		var posCfg worker.POSCfg
		if TaskName == "pos" {
			addArgs, err = worker.GetPOSArgs()
			if err != nil {
				panic(err)
			}

			_, err := posCfg.GetPOSCfg()
			if err != nil {
				panic(err)
			}
		}

		clientCfg := client.CallWorkerConfig{
			LogName:        LogName,
			TaskName:       TaskName,
			WorkerCfg:      wCfg,
			AdditionalArgs: addArgs,
			POSCfg:         posCfg,
		}

		err = job.AddJob(clientCfg)
		if err != nil {
			fmt.Printf("unable to add job: %v", err)
			return
		}
	},
}
