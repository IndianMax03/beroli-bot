package main

import (
	"context"
	"fmt"
	"sync"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var invoker *Invoker
var usersLockMap map[string]chan string = make(map[string]chan string)
var mapLock sync.Mutex

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

func sendPreliminaryMessagesWithContext(ctx context.Context, result string, err error) {
	go func() {
		messageID, err := getContextMessageID(ctx)
		if err != nil {
			panic(err)
		}
		chatID, err := getContextChatID(ctx)
		if err != nil {
			panic(err)
		}
		preliminaryMessage := PreliminaryMessage{
			chatID:    chatID,
			messageID: messageID,
			result:    result,
			err:       err,
		}
		preliminaryQueue <- preliminaryMessage
	}()
}

func preliminaryMessagesDaemon() {
	var result string
	for pm := range preliminaryQueue {
		go func() {
			if pm.err != nil {
				result = fmt.Sprintf("Ошибка: %v", pm.err)
			} else {
				result = pm.result
			}
			sendMessage(pm.chatID, pm.messageID, result)
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
		mapLock.Lock()
		if _, ok := usersLockMap[username]; !ok {
			usersLockMap[username] = make(chan string, 1)
			usersLockMap[username] <- username
		}
		mapLock.Unlock()
		<-usersLockMap[username]
		result, err = invoker.executeCommand(ctx, username, text, fileID)
		usersLockMap[username] <- username
	}
	if err != nil {
		result = fmt.Sprintf("Ошибка: %v", err)
	}
	sendMessage(update.Message.Chat.ID, update.Message.MessageID, result)
}
