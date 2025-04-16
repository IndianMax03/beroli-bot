package main

import (
	"context"
	"errors"
)

var preliminaryQueue = make(chan PreliminaryMessage, 100)

type contextKey string

var (
	ErrEmptyMessageID = errors.New("в контексте отсутствует message ID")
	ErrCastMessageID  = errors.New("message ID не int")
	ErrEmptyChatID    = errors.New("в контексте отсутствует chat ID")
	ErrCastChatID     = errors.New("chat ID не int64")
)

var messageIDKey = contextKey("messageID")
var chatIDKey = contextKey("chatID")

type PreliminaryMessage struct {
	chatID    int64
	messageID int
	result    string
	err       error
}

func getContextMessageID(ctx context.Context) (int, error) {
	messageID := ctx.Value(messageIDKey)
	if messageID == nil {
		return -1, ErrEmptyMessageID
	}
	msgID, ok := messageID.(int)
	if !ok {
		return -1, ErrCastMessageID
	}
	return msgID, nil
}

func getContextChatID(ctx context.Context) (int64, error) {
	chatID := ctx.Value(chatIDKey)
	if chatID == nil {
		return -1, ErrEmptyChatID
	}
	chtID, ok := chatID.(int64)
	if !ok {
		return -1, ErrCastChatID
	}
	return chtID, nil
}

func setContextMessageID(ctx context.Context, messageID int) (resCtx context.Context) {
	resCtx = context.WithValue(ctx, messageIDKey, messageID)
	return
}

func setContextChatID(ctx context.Context, chatID int64) (resCtx context.Context) {
	resCtx = context.WithValue(ctx, chatIDKey, chatID)
	return
}
