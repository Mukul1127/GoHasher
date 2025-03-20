package list

import (
	"github.com/Mukul1127/GoHasher/src/hashing"
	"github.com/Mukul1127/GoHasher/src/logger"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the available hashes",
	Long:  "Lists the available hashes in no specific order",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Get().Info("Available hashing algorithms")
		for _, algorithm := range hashing.GetHashFunctionNames() {
			logger.Get().Infof("- %s", algorithm)
		}
	},
}
