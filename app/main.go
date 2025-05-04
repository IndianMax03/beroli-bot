package main

import (
	"context"
	"log"

	domain "github.com/IndianMax03/beroli-bot/internal/domain"
	infra "github.com/IndianMax03/beroli-bot/internal/infra"
	telegram "github.com/IndianMax03/beroli-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := loadEnvs()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := infra.NewConnection(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := repo.CloseConnection(ctx); err != nil {
			panic(err)
		}
	}()

	initEntities(repo)

	bot := telegram.WithDebug(telegram.NewBot())

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	log.Printf("Authorized on account %s", bot.Self.UserName)

	sender := domain.NewTelegramMessageSender()

	updates := bot.GetUpdatesChan(updateConfig)

	go domain.PreliminaryMessagesDaemon(sender)

	for update := range updates {
		go domain.RunUpdate(update, sender)
	}

}
