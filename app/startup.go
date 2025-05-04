package main

import (
	"os"

	global "github.com/IndianMax03/beroli-bot/internal/global"
	domain "github.com/IndianMax03/beroli-bot/internal/domain"
	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	env "github.com/joho/godotenv"
	infra "github.com/IndianMax03/beroli-bot/internal/infra"
)

func loadEnvs() error {
	env.Load()
	global.TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	global.YANDEX_API_TOKEN = os.Getenv("YANDEX_API_TOKEN")
	global.YANDEX_ORGANIZATION_ID = os.Getenv("YANDEX_ORGANIZATION_ID")
	global.ALLOWED_USERNAME = os.Getenv("ALLOWED_USERNAME")
	global.MONGO_URL = os.Getenv("MONGO_URL")
	global.MONGO_DB_NAME = os.Getenv("MONGO_DB_NAME")
	global.MONGO_COLLECTION_NAME = os.Getenv("MONGO_COLLECTION_NAME")
	global.TRACKER_QUEUE = os.Getenv("TRACKER_QUEUE")
	global.MONGO_USER = os.Getenv("MONGO_USER")
	global.MONGO_PASSWORD = os.Getenv("MONGO_PASSWORD")
	return nil
}

func initEntities(repo *infra.MongoRepository) {
	trackerClient := api.New(global.YANDEX_API_TOKEN, global.YANDEX_ORGANIZATION_ID, "", "")
	receiver := domain.NewHandler(repo, trackerClient)
	domain.Inv = domain.NewInvoker(receiver)
}
