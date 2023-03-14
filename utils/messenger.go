package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func PrintUpdates() {
	botToken := os.Getenv("RANO_TELEGRAM_BOT_TOKEN")
	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%s", botToken)

	updatesUrl := fmt.Sprintf("%s/getUpdates", baseUrl)

	response, err := http.Get(updatesUrl)
	if err != nil {
		panic(err)
	}
	fmt.Print(response)
	defer response.Body.Close()
}

func SendToTelegram(
	chatId int,
	message string,
) error {

	botToken := os.Getenv("RANO_TELEGRAM_BOT_TOKEN")
	baseUrl := fmt.Sprintf("https://api.telegram.org/bot%s", botToken)

	url := fmt.Sprintf("%s/sendMessage", baseUrl)

	body, _ := json.Marshal(map[string]any{
		"chat_id": chatId,
		"text":    message,
	})

	_, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	return err
}
