package main

import (
	"context"
	"fmt"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
)

type Receiver interface {
	myIssues() string
	createIssue() string
	update() string
	done() string
	cancel() string
	noCommand(string) string
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
func (h Handler) myIssues() string {
	issues, _ := h.trackerClient.SearchAllIssues(
		&model.IssueSearchRequest{
			Queue: TRACKER_QUEUE,
		},
	)
	result := ""
	for i, issue := range issues {
		result += fmt.Sprintf("%v) %s\n", i, issue.Self)
	}
	return result
}

func (h Handler) createIssue() string {
	return "createIssue stub"
}

func (h Handler) update() string {
	return "update stub"
}

func (h Handler) done() string {
	return "done stub"
}

func (h Handler) cancel() string {
	return "cancel stub"
}

func (h Handler) noCommand(text string) string {
	ctx := context.Background()
	if err := h.collection.ExistsByUsername(ctx, text); err != nil {
		return fmt.Sprintf("%s, not exists. error: %v", text, err)
	}
	return fmt.Sprintf("%s, exists", text)
}
