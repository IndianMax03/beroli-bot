package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
)

type Receiver interface {
	myIssues(string) (string, error)
	createIssue(string, string) (string, error)
	done(string) (string, error)
	cancel(string) (string, error)
	noCommand(string) (string, error)
	helpCommand(map[string]Command) (string, error)
	stateCommand(string, map[string]Command) (string, error)
	ValidateState(string, string) error
	ValidateAndInitUser(string) error
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

func (h Handler) myIssues(username string) (string, error) {
	ctx := context.Background()
	issues, err := h.collection.GetIssues(ctx, username)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for i, issue := range issues {
		b.WriteString(fmt.Sprintf("%v) %s: %s\n", i+1, issue.Key, issue.Link))
	}
	result := b.String()
	if result == "" {
		result = "Вы ещё не создали ни одной задачи"
	}
	return result, nil
}

func (h Handler) createIssue(username, text string) (string, error) {
	ctx := context.Background()

	if err := h.collection.UpdateStateUser(ctx, username, CREATING_STATE); err != nil {
		return "", err
	}

	if err := h.collection.ClearIssue(ctx, username); err != nil {
		return "", err
	}

	return "Можете создавать задачу", nil
}

func (h Handler) done(username string) (string, error) {
	ctx := context.Background()

	user, err := h.collection.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	err = user.validateRequest()
	if err != nil {
		return "", err
	}

	created, err := h.trackerClient.CreateIssue(user.Issue)
	if err != nil {
		return "", err
	}

	issueData := Issue{
		Key:  created.Key,
		Link: NewIssueLink(created.Key),
	}

	err = h.collection.AppendDataIssue(ctx, username, &issueData)
	if err != nil {
		return "", err
	}

	err = h.collection.ClearIssue(ctx, username)
	if err != nil {
		return "", err
	}

	err = h.collection.UpdateStateUser(ctx, username, DONE_STATE)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Задача успешно создана: %s", issueData.Link), nil
}

func (h Handler) cancel(username string) (string, error) {
	ctx := context.Background()
	if err := h.collection.ClearIssue(ctx, username); err != nil {
		return "", err
	}
	if err := h.collection.UpdateStateUser(ctx, username, CANCELED_STATE); err != nil {
		return "", err
	}
	return "Создание задачи успешно отменено", nil
}

func (h Handler) noCommand(text string) (string, error) {
	return "no command stub", nil
}

func (h Handler) helpCommand(commandMap map[string]Command) (string, error) {
	var b strings.Builder
	for name, cmd := range commandMap {
		if name != "" {
			b.WriteString(fmt.Sprintf("%s\n", cmd.GetDescription()))
		}
	}
	result := b.String()
	if result == "" {
		result = "Ни одной команды не зарегестрировано"
	}
	return result, nil
}

func (h Handler) stateCommand(username string, commandMap map[string]Command) (string, error) {
	ctx := context.Background()
	state, err := h.collection.GetStateUser(ctx, username)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("@%s, ваш контекст: \n", username))

	stateDescr, err := GetLocalizedStateDescription(state)
	if err != nil {
		return "", err
	}
	b.WriteString(fmt.Sprintf("Состояние: '%s' \n", stateDescr))

	issues, err := h.collection.GetIssues(ctx, username)
	if err != nil {
		return "", err
	}
	b.WriteString(fmt.Sprintf("Вы создали: %v задач \n", len(issues)))

	b.WriteString("Доступные команды:\n")
	for name, cmd := range commandMap {
		if name != "" {
			if h.ValidateState(username, name) == nil {
				b.WriteString(fmt.Sprintf("- %s\n", cmd.GetDescription()))
			}
		}
	}

	if state == CREATING_STATE {
		b.WriteString("\n--- Наполнение задачи ---\n\n")
		issue, err := h.collection.GetIssue(ctx, username)
		if err != nil {
			return "", err
		}
		b.WriteString(GetLocalizedIssueFilling(issue))
	}

	return b.String(), nil
}

var (
	ErrCreatingNotStarted = errors.New("вы еще не начали создание задачи")
	ErrCreatingStarted    = errors.New("вы уже начали создание задачи")
)

func (h Handler) ValidateAndInitUser(username string) error {
	ctx := context.Background()
	if err := h.collection.ExistsUser(ctx, username); err != nil {
		user := User{
			Username: username,
			State:    NIL_STATE,
			Issues:   []Issue{},
		}
		h.collection.CreateUser(ctx, &user)
	}
	return nil
}

func (h Handler) ValidateState(username, cmdName string) error {
	ctx := context.Background()
	state, err := h.collection.GetStateUser(ctx, username)
	if err != nil {
		return err
	}

	switch state {
	case NIL_STATE, DONE_STATE, CANCELED_STATE:
		if cmdName == MY_ISSUES_COMMAND || cmdName == CREATE_ISSUE_COMMAND || cmdName == HELP_COMMAND || cmdName == STATE_COMMAND {
			return nil
		} else {
			return ErrCreatingNotStarted
		}
	case CREATING_STATE:
		if cmdName == MY_ISSUES_COMMAND || cmdName == DONE_COMMAND || cmdName == CANCEL_COMMAND || cmdName == NIL_COMMAND || cmdName == HELP_COMMAND || cmdName == STATE_COMMAND {
			return nil
		} else {
			return ErrCreatingStarted
		}
	}

	return nil
}
