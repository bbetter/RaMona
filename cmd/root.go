package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
}

var rootCmd = &cobra.Command{
	Use:   "ramona",
	Short: "RaMona is law monitor",
	Long:  `RaMona is ukrainian rada.gov.ua latest laws monitor`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
