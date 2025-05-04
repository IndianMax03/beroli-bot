package domain

import (
	telegram "github.com/IndianMax03/beroli-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMessageSender struct{}

func NewTelegramMessageSender() TelegramMessageSender {
	return TelegramMessageSender{}
}

func (TelegramMessageSender) SendMessage(chatID int64, messageID int, result string) error {
	msg := tgbotapi.NewMessage(chatID, result)
	msg.ReplyMarkup = NumericKeyboard
	msg.ReplyToMessageID = messageID
	_, err := telegram.Bot.Send(msg)
	return err
}
