package main

import (
	"context"
	"log"
	"os"

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

func loadEnvs() error {
	err := env.Load()
	if err != nil {
		return err
	}
	TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	YANDEX_API_TOKEN = os.Getenv("YANDEX_API_TOKEN")
	YANDEX_ORGANIZATION_ID = os.Getenv("YANDEX_ORGANIZATION_ID")
	ALLOWED_USERNAME = os.Getenv("ALLOWED_USERNAME")
	MONGO_URL = os.Getenv("MONGO_URL")
	MONGO_DB_NAME = os.Getenv("MONGO_DB_NAME")
	MONGO_COLLECTION_NAME = os.Getenv("MONGO_COLLECTION_NAME")
	TRACKER_QUEUE = os.Getenv("TRACKER_QUEUE")
	return nil
}

func main() {
	err := loadEnvs()
	if err != nil {
		panic(err)
	}

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

	initEntities(repo)

	bot := WithDebug(NewBot())

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.GetUpdatesChan(updateConfig)

	go delayedMessagesDaemon()

	for update := range updates {
		go runUpdate(update)
	}

}
