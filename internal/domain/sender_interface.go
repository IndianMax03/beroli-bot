package domain

type MessageSender interface {
	SendMessage(chatID int64, messageID int, result string) error
}
