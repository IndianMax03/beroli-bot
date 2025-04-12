package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewNumericKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CREATE_ISSUE_COMMAND),
			tgbotapi.NewKeyboardButton(STATE_COMMAND),
			tgbotapi.NewKeyboardButton(MY_ISSUES_COMMAND),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DONE_COMMAND),
			tgbotapi.NewKeyboardButton(HELP_COMMAND),
			tgbotapi.NewKeyboardButton(CANCEL_COMMAND),
		),
	)
}
