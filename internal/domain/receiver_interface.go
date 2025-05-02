package domain

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Receiver interface {
	MyIssues(string) (string, error)
	CreateIssue(string, string, map[string]Command) (string, error)
	Done(context.Context, string) (string, error)
	Cancel(string) (string, error)
	NoCommand(string, string, string, tgbotapi.FileID) (string, error)
	HelpCommand(map[string]Command) (string, error)
	StateCommand(string, map[string]Command) (string, error)
	ValidateState(string, string) error
	ValidateAndInitUser(string) error
	CreateTrackerIssue(dbCtx context.Context, user *User) (*Issue, error)
	UploadAttachments(user *User) error
	UploadDescriptionAttachments(user *User) error
}
