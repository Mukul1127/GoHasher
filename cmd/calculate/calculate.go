package calculate

import (
	"github.com/Mukul1127/GoHasher/src/hashing"
	"github.com/Mukul1127/GoHasher/src/logger"

	"hash"

	"github.com/spf13/cobra"
)

var CalculateCmd = &cobra.Command{
	Use:   "calculate",
	Short: "Calculates the hash of a file",
	Long:  "Calculates the hash of a file with the algorithms specified by the user in parallel",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logger.Get().Fatal("Please provide a file to be hashed")
			return
		}

		filePath := args[0]

		algorithmsFlag, err := cmd.Flags().GetStringSlice("algorithms")
		if err != nil {
			logger.Get().Fatalf("Failed to retrieve the flag 'algorithms': %s", err)
		}

		var algorithms []hash.Hash
		var algorithmLabels []string
		for _, algorithm := range algorithmsFlag {
			hashFunc, err := hashing.GetHashFunction(algorithm)
			if err != nil {
				logger.Get().Warnf("Failed to retrieve the hash function: %s", algorithm)
				continue
			}
			algorithms = append(algorithms, hashFunc)
			algorithmLabels = append(algorithmLabels, algorithm)
		}

		hashes, err := hashing.HashFile(filePath, algorithms, 1048576)
		if err != nil {
			logger.Get().Fatalf("Error hashing file '%s': %v", filePath, err)
			return
		}

		for index, hashVal := range hashes {
			logger.Get().Infof("%s: %s", algorithmLabels[index], hashVal)
		}
	},
}

func init() {
	CalculateCmd.Flags().StringSliceP("algorithms", "a", []string{}, "Comma-separated list of hash algorithms to use (e.g., md5, sha1, sha2_256)")
	err := CalculateCmd.MarkFlagRequired("algorithms")
	if err != nil {
		logger.Get().Fatalf("Failed to mark 'algorithms' flag as required: %s", err)
	}
}
