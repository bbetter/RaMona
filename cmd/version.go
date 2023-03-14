package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of RaMona",
	Long:  `All software has versions. This is RaMona's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("RaMona ver 0.0.1-alpha01")
	},
}
