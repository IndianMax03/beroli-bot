package domain

import (
	"context"
	"fmt"
	"sync"

	util "github.com/IndianMax03/beroli-bot/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Inv *Invoker
var usersLockMap map[string]chan string = make(map[string]chan string)
var mapLock sync.Mutex

func SendPreliminaryMessagesWithContext(ctx context.Context, result string, restultErr error) {
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
			err:       restultErr,
		}
		preliminaryQueue <- preliminaryMessage
	}()
}

func PreliminaryMessagesDaemon(sender MessageSender) {
	var result string
	for pm := range preliminaryQueue {
		go func() {
			if pm.err != nil {
				result = fmt.Sprintf("Ошибка: %v", pm.err)
			} else {
				result = pm.result
			}
			sender.SendMessage(pm.chatID, pm.messageID, result)
		}()
	}
}

func RunUpdate(update tgbotapi.Update, sender MessageSender) {
	if update.Message.Text == "" && update.Message.Caption == "" {
		return
	}
	ctx := context.Background()
	ctx = setContextMessageID(ctx, update.Message.MessageID)
	ctx = setContextChatID(ctx, update.Message.Chat.ID)
	var result string

	username, text, fileID, err := util.ParseMessage(update.Message)
	if err == nil {
		mapLock.Lock()
		if _, ok := usersLockMap[username]; !ok {
			usersLockMap[username] = make(chan string, 1)
			usersLockMap[username] <- username
		}
		mapLock.Unlock()
		<-usersLockMap[username]
		result, err = Inv.ExecuteCommand(ctx, username, text, fileID)
		usersLockMap[username] <- username
	}
	if err != nil {
		result = fmt.Sprintf("Ошибка: %v", err)
	}

	if err = sender.SendMessage(update.Message.Chat.ID, update.Message.MessageID, result); err != nil {
		panic(err)
	}
}
