package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
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

		bot, err := tgbotapi.NewBotAPI(os.Getenv("RAMONA_TELEGRAM_BOT_TOKEN"))
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
				result := fetchResultsAsString(cmd, cfg.filters)
				message = tgbotapi.NewMessage(input.Chat.ID, result)
				_, _ = bot.Send(message)
			case "subscribe":

				go func() {
					for {
						select {
						case <-config[input.Chat.ID].c:
							return
						default:
							var cfg = config[input.Chat.ID]
							result := fetchResultsAsString(cmd, cfg.filters)
							message := tgbotapi.NewMessage(input.Chat.ID, result)
							message.ParseMode = "HTML"
							_, _ = bot.Send(message)
							time.Sleep(60 * 60 * time.Second)
							// Do other stuff
						}
					}
				}()
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Без питань, наступний апдейт за добу",
				)

			case "unsubscribe":
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Без питань, після останнього сповіщення відпишусь",
				)
				config[input.Chat.ID].c <- true

			default:
				message = tgbotapi.NewMessage(
					input.Chat.ID,
					"Шо ти мелеш?",
				)
			}

			message.ParseMode = "HTML"
			if _, err := bot.Send(message); err != nil {
			}
		}
	},
}

func fetchResultsAsString(cmd *cobra.Command, filters []string) string {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	fetchArgs := make([]string, 2)
	fetchArgs[0] = "--filters"
	fetchArgs[1] = strings.Join(filters, ",")

	fetchCmd.Run(cmd, fetchArgs)
	return buf.String()
}
