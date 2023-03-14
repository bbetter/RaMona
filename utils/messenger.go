package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DeliveryCallback func(int) error

func SendToTelegram(
	message string,
	onMessageDelivered DeliveryCallback,
) {

	botToken := os.Getenv("RANO_TELEGRAM_BOT_TOKEN")
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

	uniqueChatIds := Distinct(Map(updates.Result, func(res TResult) int {
		return res.Message.Chat.Id
	}))

	url := fmt.Sprintf("%s/sendMessage", baseUrl)

	for _, id := range uniqueChatIds {

		body, _ := json.Marshal(map[string]any{
			"chat_id": id,
			"text":    message,
		})

		_, _ = http.Post(
			url,
			"application/json",
			bytes.NewBuffer(body),
		)

		_ = onMessageDelivered(id)
	}
}
