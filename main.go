package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
)

type Queue struct {
	Id   int
	Key  string
	Name string
}

var (
	TELEGRAM_TOKEN         string
	YANDEX_API_TOKEN       string
	YANDEX_ORGANIZATION_ID string
	ALLOWED_USERNAME       string
)

func loadEnvs() (bool, error) {
	err := env.Load()
	if err != nil {
		log.Fatal("Can't find .env file")
		return false, err
	}
	TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	YANDEX_API_TOKEN = os.Getenv("YANDEX_API_TOKEN")
	YANDEX_ORGANIZATION_ID = os.Getenv("YANDEX_ORGANIZATION_ID")
	ALLOWED_USERNAME = os.Getenv("ALLOWED_USERNAME")
	return true, nil
}

func main() {
	loadEnvs()

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_TOKEN)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	client := http.Client{Timeout: 3 * time.Second}
	uri := "https://api.tracker.yandex.net/v2/queues/"
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", YANDEX_API_TOKEN))
	req.Header.Add("X-Cloud-Org-Id", YANDEX_ORGANIZATION_ID)

	log.Printf("Authorized on account %s", bot.Self.UserName)
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil || update.Message.From.UserName != ALLOWED_USERNAME {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var queues []Queue
		json.Unmarshal(body, &queues)
		result := fmt.Sprintf("Status: %s\nQueues:\n", resp.Status)
		for _, q := range queues {
			result += fmt.Sprintf("ID = %v; Key = %s; Name = %s\n", q.Id, q.Key, q.Name)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		msg.ReplyToMessageID = update.Message.MessageID
		resp.Body.Close()
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}

}
