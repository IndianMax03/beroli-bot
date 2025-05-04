package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	global "github.com/IndianMax03/beroli-bot/internal/global"
)

const (
	TELEGRAM_BOT_API_BASE_URL = "https://api.telegram.org/file/bot"
)

var Bot *tgbotapi.BotAPI

func NewBot() *tgbotapi.BotAPI {
	var err error
	Bot, err = tgbotapi.NewBotAPI(global.TELEGRAM_TOKEN)
	if err != nil {
		panic(err)
	}
	return Bot
}

func WithDebug(bot *tgbotapi.BotAPI) *tgbotapi.BotAPI {
	bot.Debug = true
	return bot
}

func GetFileURLByPath(filePath string) (url string) {
	url = TELEGRAM_BOT_API_BASE_URL + Bot.Token + "/" + filePath
	return
}

func GetFileByConfig(fileConfig *tgbotapi.FileConfig) (tgbotapi.File, error) {
	return Bot.GetFile(*fileConfig)
}
