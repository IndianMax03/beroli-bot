package main

import (
	"errors"

	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
)

const (
	USERNAME_FIELD = "username"
	STATE_FIELD    = "state"
	ISSUE_FIELD    = "issue"
	ISSUES_FIELD   = "issues"
)

var (
	ErrMandatoryIsRequired = errors.New("mandatory fields is required")
)

type User struct {
	Username string                    `bson:"username"`
	State    string                    `bson:"state"`
	Issue    *model.IssueCreateRequest `bson:"issue"`
	Issues   []Issue                   `bson:"issues"`
}

type Issue struct {
	Key  string
	Link string
}

func NewDefaultIssue() *model.IssueCreateRequest {
	return &model.IssueCreateRequest{
		Queue: model.Queue{
			Key: TRACKER_QUEUE,
		},
	}
}

func (u *User) validateRequest() error {
	if u.Issue.Summary == "" || u.Issue.Queue.Key == "" {
		return ErrMandatoryIsRequired
	}
	return nil
}
