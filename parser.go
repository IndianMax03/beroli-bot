package main

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MAX_ATTACHMENT_SIZE_MBIT = 1024
	MAX_ATTACHMENT_SIZE_BYTE = MAX_ATTACHMENT_SIZE_MBIT * 1000 * 1000 / 8
)

var ErrAttachmentSize = errors.New("размер вложения превышает максимальный")

func ParseMessage(msg *tgbotapi.Message) (username, text string, fileID tgbotapi.FileID, err error) {
	username = msg.From.UserName
	text = msg.Text
	if len(msg.Photo) > 0 {
		for i, ph := range msg.Photo {
			if ph.FileSize > MAX_ATTACHMENT_SIZE_BYTE && i == 0 {
				return "", "", "", fmt.Errorf("%w: %v Мб", ErrAttachmentSize, MAX_ATTACHMENT_SIZE_MBIT)
			}
			if ph.FileSize <= MAX_ATTACHMENT_SIZE_BYTE {
				fileID = tgbotapi.FileID(ph.FileID)
			}
		}
	} else if msg.Document != nil {
		if msg.Document.FileSize > MAX_ATTACHMENT_SIZE_BYTE {
			return "", "", "", fmt.Errorf("%w: %v Мб", ErrAttachmentSize, MAX_ATTACHMENT_SIZE_MBIT)
		}
		fileID = tgbotapi.FileID(msg.Document.FileID)
	}
	if fileID != "" {
		text = msg.Caption
	}
	return
}
