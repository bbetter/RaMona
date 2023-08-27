package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
	"owl.com/ramona/utils"
)

type UserConfig struct {
	filters []string
	c       chan bool
}

var config map[int64]UserConfig

func init() {

	config = make(map[int64]UserConfig, 10)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "telegram bot subscriptions monitor",
	Long:  `telegram bot subscriptions monitor`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: subscribe to updates channel , get chat id
		// register periodic job for that specific chat id

		execPath, _ := os.Executable()
		logFilePath := fmt.Sprintf("%s\\logs.txt", filepath.Dir(execPath))

		//setup logging
		f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open logs file")
		}
		log.SetOutput(f)
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		bot, err := tgbotapi.NewBotAPI(os.Getenv("RANO_TELEGRAM_BOT_TOKEN"))
		if err != nil {
			panic(err)
		}

		bot.Debug = true
		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 60
		updates := bot.GetUpdatesChan(updateConfig)
		for update := range updates {
			if update.Message == nil {
				continue
			}

			if !update.Message.IsCommand() {
				continue
			}

			var input = update.Message

			var message tgbotapi.MessageConfig
			switch input.Command() {
			case "start":
				config[input.Chat.ID] = UserConfig{
					[]string{},
					make(chan bool),
				}
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Привіт. Я бот помічник - Рамона. Ось команди які я підтримую (filters, schedule, setfilters, setschedule, fetch subscribe)",
				)
			case "filters":
				var cfg = config[input.Chat.ID]
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					fmt.Sprintf("Наразі встановлені такі фільтри: [%s]", strings.Join(cfg.filters, " ")),
				)
			case "setfilters":
				var cfg = config[input.Chat.ID]
				config[input.Chat.ID] = UserConfig{
					strings.Split(input.CommandArguments(), " "),
					cfg.c,
				}
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					fmt.Sprintf("Готово: %s", strings.Join(config[input.Chat.ID].filters, " ")),
				)
			case "fetch":
				var cfg = config[input.Chat.ID]
				message = prepareMessage(
					cfg.filters,
					input.Chat.ID,
				)
			case "subscribe":

				go func() {
					for {
						select {
						case <-config[input.Chat.ID].c:
							return
						default:
							var cfg = config[input.Chat.ID]
							message = prepareMessage(
								cfg.filters,
								input.Chat.ID,
							)
							_, _ = bot.Send(message)
							time.Sleep(60 * 60 * time.Second)
							// Do other stuff
						}
					}
				}()
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Без питань - наступний апдейт за добу",
				)

			case "unsubscribe":
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Без питань - після останнього сповіщення відпишусь",
				)
				config[input.Chat.ID].c <- true

			default:
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Шо ти мелеш?",
				)
			}

			if _, err := bot.Send(message); err != nil {
			}
		}
	},
}

func prepareMessage(filters []string, chatId int64) tgbotapi.MessageConfig {
	utils.TimeLog("\n Завантаження даних.")

	items := utils.ParseFeedItems()
	utils.TimeLog(fmt.Sprintf("Дані завантажено. Загальна к-сть: %d", len(items)))

	if len(filters) != 0 {
		utils.TimeLog(fmt.Sprintf("Пошук...( %v )", filters))
		items = utils.FilterByTriggers(items, filters)
		utils.TimeLog(fmt.Sprintf("Пошук завершено. К-сть співпадінь: %d", len(items)))
	}

	if len(items) == 0 {
		msg := tgbotapi.NewMessage(chatId, "Нажаль я <b><u>НІЧОГО</u></b> не знайшов спробуй інші слова. :()")
		msg.ParseMode = "HTML"
		return msg
	}

	messages := utils.Map(items, func(item *gofeed.Item) string {
		return fmt.Sprintf("%s\n%s", item.Description, item.Link)
	})

	message := strings.Join(messages, "\n\n")

	//highlight occurences
	var fRegexp *regexp.Regexp
	for _, filter := range filters {
		fRegexp = regexp.MustCompile(fmt.Sprintf(`(?i)%s`, filter))
		message = fRegexp.ReplaceAllString(message, fmt.Sprintf("<b><u>%s</u></b>", strings.ToUpper(filter)))
	}

	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = "HTML"
	return msg
}
