package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ramona",
	Short: "RaMona is law monitor",
	Long:  `RaMona is ukrainian rada.gov.ua latest laws monitor`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(versionCmd)
}
