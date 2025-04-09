package main

import (
	"context"
	"errors"
	"fmt"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
)

type Receiver interface {
	myIssues() (string, error)
	createIssue(string, string) (string, error)
	update() (string, error)
	done() (string, error)
	cancel() (string, error)
	noCommand(string) (string, error)
	ValidateState(string, string) error
}

type Handler struct {
	collection    *Collection
	trackerClient *api.Client
}

func NewHandler(repo *MongoRepository, trackerClient *api.Client) Receiver {
	return Handler{
		collection:    NewCollection(*repo),
		trackerClient: trackerClient,
	}
}

// TODO: temporary solution to check out tracker integration
func (h Handler) myIssues() (string, error) {
	issues, _ := h.trackerClient.SearchAllIssues(
		&model.IssueSearchRequest{
			Queue: TRACKER_QUEUE,
		},
	)
	result := ""
	for i, issue := range issues {
		result += fmt.Sprintf("%v) %s\n", i, issue.Self)
	}
	return result, nil
}

func (h Handler) createIssue(username, text string) (string, error) {
	ctx := context.Background()
	if err := h.collection.ExistsUser(ctx, username); err != nil {
		user := User{
			Username: username,
			State:    CREATING_STATE,
		}
		h.collection.CreateUser(ctx, &user)
	} else {
		h.collection.UpdateStateUser(ctx, username, CREATING_STATE)
	}
	return "ok", nil
}

func (h Handler) update() (string, error) {
	return "update stub", nil
}

func (h Handler) done() (string, error) {
	return "done stub", nil
}

func (h Handler) cancel() (string, error) {
	return "cancel stub", nil
}

func (h Handler) noCommand(text string) (string, error) {
	ctx := context.Background()
	if err := h.collection.ExistsUser(ctx, text); err != nil {
		return fmt.Sprintf("%s, not exists. error: %v", text, err), nil
	}
	return fmt.Sprintf("%s, exists", text), nil
}

var (
	ErrCreatingNotStarted = errors.New("creating not started")
	ErrCreatingStarted    = errors.New("creating started")
	ErrCreatingInError    = errors.New("creating in error")
)

func (h Handler) ValidateState(username, cmdName string) error {
	ctx := context.Background()
	state, err := h.collection.GetStateUser(ctx, username)
	if err != nil {
		return err
	}

	switch state {
	case NIL_STATE, DONE_STATE, CANCELED_STATE:
		if cmdName == MY_ISSUES_COMMAND || cmdName == CREATE_ISSUE_COMMAND {
			return nil
		} else {
			return ErrCreatingNotStarted
		}
	case CREATING_STATE:
		if cmdName == MY_ISSUES_COMMAND || cmdName == DONE_COMMAND || cmdName == CANCEL_COMMAND || cmdName == NIL_COMMAND {
			return nil
		} else {
			return ErrCreatingStarted
		}
	case ERROR_CREATING_STATE:
		if cmdName == MY_ISSUES_COMMAND || cmdName == CANCEL_COMMAND || cmdName == NIL_COMMAND {
			return nil
		} else {
			return ErrCreatingInError
		}
	}

	return nil
}
