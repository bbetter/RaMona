package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/context"
)

const feedUrl = "https://feeds.feedburner.com/gov/gnjU"

func main() {

	execPath, _ := os.Executable()
	logFilePath := fmt.Sprintf("%s\\logs.txt", filepath.Dir(execPath))

	//setup logging
	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open logs file")
	}
	defer f.Close()
	log.SetOutput(f)

	// read variables
	botToken := readEnvVars()
	triggers := readFlags()

	logWithCurrentTime("\n Завантаження даних.")

	allFeedItems := parseFeedItems()
	logWithCurrentTime(fmt.Sprintf("Дані завантажено. Загальна к-сть: %d", len(allFeedItems)))

	logWithCurrentTime(fmt.Sprintf("Пошук...( %v )", triggers))
	filteredFeedItems := filterByTriggers(allFeedItems, triggers)
	logWithCurrentTime(fmt.Sprintf("Пошук завершено. К-сть співпадінь: %d", len(filteredFeedItems)))

	if len(filteredFeedItems) == 0 {
		return
	}

	messages := Map(filteredFeedItems, func(item *gofeed.Item) string {
		return fmt.Sprintf("%s\n%s", item.Description, item.Link)
	})
	message := strings.Join(messages, "\n\n")

	if len(botToken) == 0 {
		logWithCurrentTime("Відсутні змінні середовища для налаштування сповіщень.")
		return
	}

	sendToTelegram(botToken, message)
}

func readEnvVars() (botToken string) {

	botToken = os.Getenv("RANO_TELEGRAM_BOT_TOKEN")

	return
}

func readFlags() (triggers []string) {
	triggersStr := flag.String("triggers", "", "список слів як можна використовувати для пошуку")
	flag.Parse()
	triggers = strings.Split(*triggersStr, " ")

	return
}

func parseFeedItems() []*gofeed.Item {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
	return feed.Items
}

func filterByTriggers(items []*gofeed.Item, triggers []string) []*gofeed.Item {

	return Filter(items, func(item *gofeed.Item) bool {
		//equalfold not working!!

		titleLowercase := strings.ToLower(item.Title)
		descLowercase := strings.ToLower(item.Description)

		byTitle := Any(triggers, func(s string) bool {
			return strings.Contains(titleLowercase, s)
		})

		byDescription := Any(triggers, func(s string) bool {
			return strings.Contains(descLowercase, s)
		})

		return byTitle || byDescription
	})
}

func sendToTelegram(botToken string, message string) {
	var execPath, _ = os.Executable()
	var filePath = fmt.Sprintf("%s\\users.csv", filepath.Dir(execPath))
	//setup database
	var f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open logs file")
	}
	defer f.Close()

	reader := csv.NewReader(f)
	writer := csv.NewWriter(f)

	csvChatIdsStr, _ := reader.Read()
	csvChatIds := Map(csvChatIdsStr, func(r string) int {
		res, _ := strconv.Atoi(r)
		return res
	})

	logWithCurrentTime(fmt.Sprintf("К-сть підписників у файлі: %d", len(csvChatIds)))
	logWithCurrentTime("Оновлюю підписників...")

	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%s", botToken)

	updatesUrl := fmt.Sprintf("%s/getUpdates", baseUrl)

	response, err := http.Get(updatesUrl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Parse the updates JSON
	var updates TUpdate
	err = json.NewDecoder(response.Body).Decode(&updates)
	if err != nil {
		panic(err)
	}

	netChatIds := Map(updates.Result, func(res TResult) int {
		return res.Message.Chat.Id
	})

	resultChatIds := Filter(Distinct(append(csvChatIds, netChatIds...)), func(n int) bool {
		return n != 0
	})

	writer.Write(Map(resultChatIds, func(s int) string {
		return strconv.Itoa(s)
	}))
	defer writer.Flush()

	url := fmt.Sprintf("%s/sendMessage", baseUrl)

	for _, id := range resultChatIds {

		body, _ := json.Marshal(map[string]any{
			"chat_id": id,
			"text":    message,
		})

		_, _ = http.Post(
			url,
			"application/json",
			bytes.NewBuffer(body),
		)

		logWithCurrentTime(fmt.Sprintf("Сповіщення доставлено до %d", id))
	}
}

func logWithCurrentTime(message string) {
	t := time.Now()
	formattedTime := t.Format("02.01.2006 15:04")

	fmt.Printf("%s# %s \n", formattedTime, message)
	log.Printf("%s# %s \n", formattedTime, message)
}
