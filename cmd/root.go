package cmd

import (
	"github.com/Mukul1127/GoHasher/cmd/calculate"
	"github.com/Mukul1127/GoHasher/cmd/list"
	"github.com/Mukul1127/GoHasher/src/logger"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "GoHasher",
	Short: "",
	Long:  "A simple, quick file hasher built in Golang!",
}

func Execute() {
	logger.Initialize()

	RootCmd.AddCommand(calculate.CalculateCmd)
	RootCmd.AddCommand(list.ListCmd)
	err := RootCmd.Execute()
	if err != nil {
		logger.Get().Fatalf("Failed to initialize Cobra: %s", err)
	}
}
