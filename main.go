package main

import (
	"context"
	"log"
	"os"
	"strings"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
)

var (
	TELEGRAM_TOKEN         string
	YANDEX_API_TOKEN       string
	YANDEX_ORGANIZATION_ID string
	ALLOWED_USERNAME       string
	MONGO_URL              string
	MONGO_DB_NAME          string
	MONGO_COLLECTION_NAME  string
	TRACKER_QUEUE          string
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
	MONGO_URL = os.Getenv("MONGO_URL")
	MONGO_DB_NAME = os.Getenv("MONGO_DB_NAME")
	MONGO_COLLECTION_NAME = os.Getenv("MONGO_COLLECTION_NAME")
	TRACKER_QUEUE = os.Getenv("TRACKER_QUEUE")
	return true, nil
}

func main() {
	loadEnvs()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := NewConnection(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := repo.CloseConnection(ctx); err != nil {
			panic(err)
		}
	}()

	trackerClient := api.New(YANDEX_API_TOKEN, YANDEX_ORGANIZATION_ID, "", "")

	receiver := NewHandler(repo, trackerClient)
	invoker := NewInvoker(&receiver)

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_TOKEN)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil || update.Message.Text == "" || update.Message.From.UserName != ALLOWED_USERNAME {
			continue
		}
		result := ""

		text := strings.Split(update.Message.Text, " ")
		if len(text) == 1 {
			result = invoker.executeCommand(text[0], "")
		} else if len(text) == 2 {
			result = invoker.executeCommand(text[0], text[1])
		} else {
			result = "deny"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}

}
