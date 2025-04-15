package main

import (
	"context"
	"fmt"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var invoker *Invoker

func initEntities(repo *MongoRepository) {
	trackerClient := api.New(YANDEX_API_TOKEN, YANDEX_ORGANIZATION_ID, "", "")
	receiver := NewHandler(repo, trackerClient)
	invoker = NewInvoker(receiver)
}

func sendMessage(chatID int64, messageID int, result string) {
	msg := tgbotapi.NewMessage(chatID, result)
	msg.ReplyMarkup = NumericKeyboard
	msg.ReplyToMessageID = messageID
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func delayedMessagesDaemon() {
	var result string
	for dm := range delayQueue {
		go func() {
			if dm.err != nil {
				result = fmt.Sprintf("Ошибка: %v", dm.err)
			} else {
				result = dm.result
			}
			sendMessage(dm.chatID, dm.messageID, result)
		}()
	}
}

func runUpdate(update tgbotapi.Update) {
	if update.Message.Text == "" && update.Message.Caption == "" {
		return
	}
	ctx := context.Background()
	ctx = setContextMessageID(ctx, update.Message.MessageID)
	ctx = setContextChatID(ctx, update.Message.Chat.ID)
	var result string

	username, text, fileID, err := ParseMessage(update.Message)
	if err == nil {
		result, err = invoker.executeCommand(ctx, username, text, fileID)
	}
	if err != nil {
		result = fmt.Sprintf("Ошибка: %v", err)
	}
	sendMessage(update.Message.Chat.ID, update.Message.MessageID, result)
}
