package main

import (
	"owl.com/ramona/cmd"
)

const feedUrl = "https://feeds.feedburner.com/gov/gnjU"

func main() {

	cmd.Execute()

	// // read variables
	// botToken := readEnvVars()
	// triggers := readFlags()

	// logWithCurrentTime("\n Завантаження даних.")

	// allFeedItems := parseFeedItems()
	// logWithCurrentTime(fmt.Sprintf("Дані завантажено. Загальна к-сть: %d", len(allFeedItems)))

	// logWithCurrentTime(fmt.Sprintf("Пошук...( %v )", triggers))
	// filteredFeedItems := filterByTriggers(allFeedItems, triggers)
	// logWithCurrentTime(fmt.Sprintf("Пошук завершено. К-сть співпадінь: %d", len(filteredFeedItems)))

	// if len(filteredFeedItems) == 0 {
	// 	return
	// }

	// messages := Map(filteredFeedItems, func(item *gofeed.Item) string {
	// 	return fmt.Sprintf("%s\n%s", item.Description, item.Link)
	// })
	// message := strings.Join(messages, "\n\n")
	// fmt.Print(message)

	// if len(botToken) == 0 {
	// 	logWithCurrentTime("Відсутні змінні середовища для налаштування сповіщень.")
	// 	return
	// }

	// sendToTelegram(botToken, message)
}
