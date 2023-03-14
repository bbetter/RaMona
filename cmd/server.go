package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "telegram bot subscriptions monitor",
	Long:  `telegram bot subscriptions monitor`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: subscribe to updates channel , get chat id
		// register periodic job for that specific chat id
	},
}
