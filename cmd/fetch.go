package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"owl.com/ramona/utils"
)

var filters []string
var chats []int

func init() {
	fetchCmd.Flags().StringSliceVarP(&filters, "filters", "f", []string{""}, "space separated triggers")
	fetchCmd.Flags().IntSliceVarP(&chats, "chats", "c", []int{}, "send to telegram bot")

	viper.AddConfigPath("./.configs")
	viper.SetConfigName("tg_bot_config")

	viper.BindPFlags(fetchCmd.Flags())

	viper.ReadInConfig()

}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch new laws",
	Long:  `Fetch latest incoming laws`,
	Run: func(cmd *cobra.Command, args []string) {

		f_chats := viper.GetIntSlice("chats")
		f_filters := viper.GetStringSlice("filters")

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

		if len(filters) != 0 {
			utils.TimeLog(fmt.Sprintf("Пошук...( %v )", f_filters))
			items = utils.FilterByTriggers(items, f_filters)
			utils.TimeLog(fmt.Sprintf("Пошук завершено. К-сть співпадінь: %d", len(items)))
		}

		messages := utils.Map(items, func(item *gofeed.Item) string {
			return fmt.Sprintf("%s\n%s", item.Description, item.Link)
		})
		message := strings.Join(messages, "\n\n")

		//highlight occurences
		for _, filter := range f_filters {
			message = strings.ReplaceAll(message, filter, fmt.Sprintf("*%s*", filter))
		}

		utils.TimeLog(message)

		for _, chatId := range f_chats {

			err := utils.SendToTelegram(chatId, message)
			if err == nil {
				utils.TimeLog(fmt.Sprintf("Повідомлення до чату %d успішно доставлено", chatId))
			}
		}
	},
}
