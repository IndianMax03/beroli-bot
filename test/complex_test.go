package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/IndianMax03/beroli-bot/internal/domain"
	"github.com/IndianMax03/beroli-bot/internal/global"
	"github.com/IndianMax03/beroli-bot/internal/infra"
	"github.com/IndianMax03/beroli-bot/internal/telegram"
	"github.com/IndianMax03/beroli-bot/mocks"
	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
	"go.uber.org/mock/gomock"
)

var (
	ctx        context.Context
	cancel     context.CancelFunc
	mockSender *mocks.MockMessageSender
	repo       *infra.MongoRepository
	ctrl       *gomock.Controller
)

var ATTACHMENT_SAMPLE_2MB string

const SIZE_2MB = 1871957

func generateCommandFromDifferentUsersChannel(count int, text string) tgbotapi.UpdatesChannel {
	channel := make(chan tgbotapi.Update, count)
	defer close(channel)

	for i := range count {
		channel <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Text: text,
			},
		}
	}
	return channel
}

func generateNilCommandFromDifferentUsersWithTextChannel(count int, textLength int) tgbotapi.UpdatesChannel {
	channel := make(chan tgbotapi.Update, count)
	defer close(channel)

	for i := range count {
		channel <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Text: fmt.Sprintf("%s %s", domain.ISSUE_SUMMARY_TAG, RandStringWithSymbolsAndEmojis(textLength)),
			},
		}
	}
	return channel
}

func generateNilCommandFromDifferentUsersWithAttachmentChannel(count int) tgbotapi.UpdatesChannel {
	channel := make(chan tgbotapi.Update, count)
	defer close(channel)

	for i := range count {
		channel <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Caption: domain.ISSUE_ATTACHMENT_TAG,
				Photo: []tgbotapi.PhotoSize{
					{
						FileID:   ATTACHMENT_SAMPLE_2MB,
						FileSize: SIZE_2MB,
					},
				},
			},
		}
	}
	return channel
}

func generatePrepareToDoneCommandFromDifferentUsersWithAttachmentChannel(count int) (tgbotapi.UpdatesChannel, tgbotapi.UpdatesChannel, tgbotapi.UpdatesChannel, tgbotapi.UpdatesChannel) {
	cancelChan := make(chan tgbotapi.Update, count)
	issueChan := make(chan tgbotapi.Update, count)
	sumChan := make(chan tgbotapi.Update, count)
	attChan := make(chan tgbotapi.Update, count)
	defer close(cancelChan)
	defer close(issueChan)
	defer close(sumChan)
	defer close(attChan)

	for i := range count {
		cancelChan <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Text: domain.CANCEL_COMMAND,
			},
		}

		issueChan <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Text: domain.CREATE_ISSUE_COMMAND,
			},
		}

		sumChan <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Text: fmt.Sprintf("[%s-%d] %s %s", usernameStub, i, domain.ISSUE_SUMMARY_TAG, RandStringWithSymbolsAndEmojis(25)),
			},
		}

		attChan <- tgbotapi.Update{
			Message: &tgbotapi.Message{
				MessageID: i,
				Chat: &tgbotapi.Chat{
					ID: int64(i),
				},
				From: &tgbotapi.User{
					ID:       int64(i),
					UserName: fmt.Sprintf("%s-%s", usernameStub, fmt.Sprint(i)),
				},
				Caption: domain.ISSUE_ATTACHMENT_TAG,
				Photo: []tgbotapi.PhotoSize{
					{
						FileID:   ATTACHMENT_SAMPLE_2MB,
						FileSize: SIZE_2MB,
					},
				},
			},
		}
	}
	return cancelChan, issueChan, sumChan, attChan
}

func initBSession(b *testing.B) {
	err := env.Load()
	if err != nil {
		panic(err)
	}
	global.TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	global.YANDEX_API_TOKEN = os.Getenv("YANDEX_API_TOKEN")
	global.YANDEX_ORGANIZATION_ID = os.Getenv("YANDEX_ORGANIZATION_ID")
	global.ALLOWED_USERNAME = os.Getenv("ALLOWED_USERNAME")
	global.MONGO_URL = os.Getenv("MONGO_URL")
	global.MONGO_DB_NAME = os.Getenv("MONGO_DB_NAME")
	global.MONGO_COLLECTION_NAME = os.Getenv("MONGO_COLLECTION_NAME")
	global.TRACKER_QUEUE = os.Getenv("TRACKER_QUEUE")
	global.MONGO_USER = os.Getenv("MONGO_USER")
	global.MONGO_PASSWORD = os.Getenv("MONGO_PASSWORD")
	ATTACHMENT_SAMPLE_2MB = os.Getenv("ATTACHMENT_SAMPLE_2MB")

	ctrl = gomock.NewController(b)
	mockSender = mocks.NewMockMessageSender(ctrl)
	ctx, cancel = context.WithCancel(context.Background())

	repo, err = infra.NewConnection(ctx)
	if err != nil {
		panic(err)
	}

	trackerClient := api.New(global.YANDEX_API_TOKEN, global.YANDEX_ORGANIZATION_ID, "", "")
	receiver := domain.NewHandler(repo, trackerClient)
	domain.Inv = domain.NewInvoker(receiver)
}

func closeBSession() {
	if err := repo.CloseConnection(ctx); err != nil {
		panic(err)
	}
	cancel()
	ctrl.Finish()
}

func BenchmarkStateCommand(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("StateCommand-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateCommandFromDifferentUsersChannel(userCount[i], domain.STATE_COMMAND)
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkHelpCommand(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("HelpCommand-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateCommandFromDifferentUsersChannel(userCount[i], domain.HELP_COMMAND)
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkCancelCommand(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("CancelCommand-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateCommandFromDifferentUsersChannel(userCount[i], domain.CANCEL_COMMAND)
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkMyIssuesCommand(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("MyIssuesCommand-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateCommandFromDifferentUsersChannel(userCount[i], domain.MY_ISSUES_COMMAND)
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkCreateIssueCommand(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("CreateIssueCommand-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateCommandFromDifferentUsersChannel(userCount[i], domain.CREATE_ISSUE_COMMAND)
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkNilCommandWithText(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("NilCommandWithText-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateNilCommandFromDifferentUsersWithTextChannel(userCount[i], 128)
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkNilCommandWithAttachment(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("NilCommandWithAttachment-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				updates := generateNilCommandFromDifferentUsersWithAttachmentChannel(userCount[i])
				done := make(chan struct{}, userCount[i])

				mockSender.
					EXPECT().
					SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).
					Times(userCount[i])

				go domain.PreliminaryMessagesDaemon(mockSender)

				for update := range updates {
					go func(u tgbotapi.Update) {
						domain.RunUpdate(u, mockSender)
						done <- struct{}{}
					}(update)
				}

				for range userCount[i] {
					<-done
				}
			}
		})
	}
}

func BenchmarkDoneCommand(b *testing.B) {
	userCount := []int{1, 20, 50, 100, 150, 300, 500, 800}

	initBSession(b)
	defer closeBSession()

	telegram.NewBot()

	b.ResetTimer()

	for i := range userCount {
		name := fmt.Sprintf("NilCommandWithAttachment-%d-users", userCount[i])
		b.Run(name, func(b *testing.B) {
			cC, iC, sC, aC := generatePrepareToDoneCommandFromDifferentUsersWithAttachmentChannel(userCount[i])
			cD := make(chan struct{}, userCount[i])
			iD := make(chan struct{}, userCount[i])
			sD := make(chan struct{}, userCount[i])
			aD := make(chan struct{}, userCount[i])

			mockSender.
				EXPECT().
				SendMessage(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).
				AnyTimes()

			go domain.PreliminaryMessagesDaemon(mockSender)

			for update := range cC {
				go func(u tgbotapi.Update) {
					domain.RunUpdate(u, mockSender)
					cD <- struct{}{}
				}(update)
			}

			for range userCount[i] {
				<-cD
			}

			for update := range iC {
				go func(u tgbotapi.Update) {
					domain.RunUpdate(u, mockSender)
					iD <- struct{}{}
				}(update)
			}

			for range userCount[i] {
				<-iD
			}

			for update := range sC {
				go func(u tgbotapi.Update) {
					domain.RunUpdate(u, mockSender)
					sD <- struct{}{}
				}(update)
			}

			for range userCount[i] {
				<-sD
			}

			for update := range aC {
				go func(u tgbotapi.Update) {
					domain.RunUpdate(u, mockSender)
					aD <- struct{}{}
				}(update)
			}

			for range userCount[i] {
				<-aD
			}

			b.ResetTimer()

			updates := generateCommandFromDifferentUsersChannel(userCount[i], domain.DONE_COMMAND)
			done := make(chan struct{}, userCount[i])

			for update := range updates {
				go func(u tgbotapi.Update) {
					domain.RunUpdate(u, mockSender)
					done <- struct{}{}
				}(update)
			}

			for range userCount[i] {
				<-done
			}

		})
	}
}
