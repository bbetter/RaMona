package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/mmcdole/gofeed"
    "golang.org/x/net/context"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)
const feedUrl = "https://feeds.feedburner.com/gov/gnjU"


func main() {

    execPath, _ := os.Executable()
    fmt.Println(execPath)
    logFilePath := fmt.Sprintf("%s\\logs.txt",filepath.Dir(execPath))
    fmt.Println(logFilePath)

    //setup logging
    f, err := os.OpenFile(logFilePath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open logs file")
    }
    defer f.Close()
    log.SetOutput(f)

    // read variables
    botToken, chatId := readEnvVars()
    triggers := readParams()

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
    message := strings.Join(messages,"\n\n")
    fmt.Print(message)

    if len(botToken) == 0 || len(chatId) == 0 {
        logWithCurrentTime("Відсутні змінні середовища для налаштування сповіщень.")
        return
    }

    sendToTelegram(botToken, chatId, message)
}

func readEnvVars() (botToken string, chatId string){

    botToken = os.Getenv("RANO_TELEGRAM_BOT_TOKEN")
    chatId = os.Getenv("RANO_TELEGRAM_CHAT_ID")

    return
}

func readParams() (triggers []string){
    triggersStr := flag.String("triggers", "", "список слів як можна використовувати для пошуку")
    flag.Parse()
    triggers = strings.Split(*triggersStr, " ")

    return
}

func parseFeedItems() [] *gofeed.Item {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    fp := gofeed.NewParser()
    feed, _ := fp.ParseURLWithContext(feedUrl, ctx)
    return feed.Items
}

func filterByTriggers(items [] *gofeed.Item, triggers [] string) [] *gofeed.Item{


    return Filter(items, func(item *gofeed.Item) bool {
        //equalfold not working!!

        titleLowercase :=strings.ToLower(item.Title)
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

func sendToTelegram(botToken string, chatId string,  message string)  {

    baseUrl := fmt.Sprintf("https://api.telegram.org/bot%s", botToken)
    url := fmt.Sprintf("%s/sendMessage", baseUrl)

    body, _ := json.Marshal(map[string]string{
        "chat_id": chatId,
        "text":    message,
        })

    _, _ = http.Post(
        url,
        "application/json",
        bytes.NewBuffer(body),
        )

    logWithCurrentTime("Сповіщення доставлено.")
}

func logWithCurrentTime(message string){
    t := time.Now()
    formattedTime := t.Format("02.01.2006 15:04")

    fmt.Printf("%s# %s \n", formattedTime, message)
    log.Printf("%s# %s \n", formattedTime, message)
}