package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
	"owl.com/ramona/utils"
)

var triggers []string
var telegram bool = false

func init() {
	fetchCmd.Flags().StringArrayVar(&triggers, "tgs", []string{""}, "space separated triggers")
	fetchCmd.Flags().BoolVar(&telegram, "tel", false, "send to telegram bot")
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch new laws",
	Long:  `Fetch latest incoming laws`,
	Run: func(cmd *cobra.Command, args []string) {

		execPath, _ := os.Executable()
		logFilePath := fmt.Sprintf("%s\\logs.txt", filepath.Dir(execPath))

		//setup logging
		f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open logs file")
		}
		defer f.Close()
		log.SetOutput(f)

		utils.TimeLog("\n Завантаження даних.")

		items := utils.ParseFeedItems()
		utils.TimeLog(fmt.Sprintf("Дані завантажено. Загальна к-сть: %d", len(items)))

		if len(triggers) != 0 {
			utils.TimeLog(fmt.Sprintf("Пошук...( %v )", triggers))
			items = utils.FilterByTriggers(items, triggers)
			utils.TimeLog(fmt.Sprintf("Пошук завершено. К-сть співпадінь: %d", len(items)))
		}

		messages := utils.Map(items, func(item *gofeed.Item) string {
			return fmt.Sprintf("%s\n%s", item.Description, item.Link)
		})
		message := strings.Join(messages, "\n\n")

		utils.TimeLog(message)

		if telegram {
			utils.SendToTelegram(message, func(i int) error {
				utils.TimeLog(fmt.Sprintf("Message to chat %d delivered successfully", i))
				return nil
			})
		}

	},
}
