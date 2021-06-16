package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/watercompany/multiplex/job"
	"github.com/watercompany/multiplex/worker"
	"github.com/watercompany/multiplex/worker/client"
)

var (
	LogName      string
	TaskName     string
	databaseAddr string
	localDrive   string
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

	addJobCmd.Flags().StringVar(
		&databaseAddr,
		"db-addr",
		"",
		"address of database to save jobs",
	)

	addJobCmd.Flags().StringVar(
		&localDrive,
		"local-drive",
		"",
		"name of local drive for temp plot files",
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

		var tempDir string = ""
		var temp2Dir string = ""
		var finalDir string = ""
		if localDrive != "" {
			tempDir = fmt.Sprintf("/mnt/%s/plotfiles/temp", localDrive)
			temp2Dir = fmt.Sprintf("/mnt/%s/plotfiles/temp2", localDrive)
			finalDir = fmt.Sprintf("/mnt/%s/plotfiles/final", localDrive)
		}

		var addArgs []string
		var posCfg worker.POSCfg
		switch TaskName {
		case "pos":
			wCfg.ExecDir = "./execs/ProofOfSpace"
			addArgs, posCfg, err = worker.GetPOSArgs(tempDir, temp2Dir, finalDir)
			if err != nil {
				panic(err)
			}
		case "posv2":
			wCfg.ExecDir = "./execs/chia_plot"
			addArgs, posCfg, err = worker.GetPOSArgs_V2(tempDir, temp2Dir, finalDir)
			if err != nil {
				panic(err)
			}
		}

		if tempDir != "" {
			posCfg.TempDir = tempDir
			posCfg.Temp2Dir = temp2Dir
			posCfg.FinalDir = finalDir
		}

		clientCfg := client.CallWorkerConfig{
			LogName:        LogName,
			TaskName:       TaskName,
			WorkerCfg:      wCfg,
			AdditionalArgs: addArgs,
			POSCfg:         posCfg,
		}

		err = job.AddJob(clientCfg, databaseAddr)
		if err != nil {
			fmt.Printf("unable to add job: %v", err)
			return
		}
	},
}
