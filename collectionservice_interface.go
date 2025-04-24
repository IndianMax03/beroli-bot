package main

import (
	"context"

	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
)

type CollectionService interface {
	CreateUser(ctx context.Context, u *User)
	UpdateUser(ctx context.Context, u *User) error
	ResetUser(ctx context.Context, username string) error
	GetUser(ctx context.Context, username string) (*User, error)
	ExistsUser(ctx context.Context, username string) error
	UpdateStateUser(ctx context.Context, username, state string) error
	GetStateUser(ctx context.Context, username string) (string, error)
	ClearIssue(ctx context.Context, username string) error
	AppendDataIssue(ctx context.Context, username string, data *Issue) error
	GetIssues(ctx context.Context, username string) ([]Issue, error)
	GetIssue(ctx context.Context, username string) (*model.IssueCreateRequest, error)
	UpdateSummaryIssue(ctx context.Context, username, text string) error
	UpdateDescriptionIssue(ctx context.Context, username, text string) error
	AppendAttachmentIssue(ctx context.Context, username string, attachmentID string) error
	AppendDescriptionAttachmentIssue(ctx context.Context, username string, descriptionAttachmentID string) error
	AppendTagIssue(ctx context.Context, username string, tags []string) error
}
