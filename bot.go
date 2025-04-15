package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	TELEGRAM_BOT_API_BASE_URL = "https://api.telegram.org/file/bot"
)

var bot *tgbotapi.BotAPI

func NewBot() *tgbotapi.BotAPI {
	var err error
	bot, err = tgbotapi.NewBotAPI(TELEGRAM_TOKEN)
	if err != nil {
		panic(err)
	}
	return bot
}

func WithDebug(bot *tgbotapi.BotAPI) *tgbotapi.BotAPI {
	bot.Debug = true
	return bot
}

func getFileURLByPath(filePath string) (url string) {
	url = TELEGRAM_BOT_API_BASE_URL + bot.Token + "/" + filePath
	return
}

func getFileByConfig(fileConfig *tgbotapi.FileConfig) (tgbotapi.File, error) {
	return bot.GetFile(*fileConfig)
}
