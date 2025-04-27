package domain

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	api "github.com/IndianMax03/yandex-tracker-go-client/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	infra "github.com/IndianMax03/beroli-bot/internal/infra"
	telegram "github.com/IndianMax03/beroli-bot/internal/telegram"
	"resty.dev/v3"
)

var (
	ErrEmptyTag        = errors.New("не указан тэг")
	ErrEmptyAttachment = errors.New("вложение не представлено")
)

type Receiver interface {
	myIssues(string) (string, error)
	createIssue(string, string, map[string]Command) (string, error)
	done(context.Context, string) (string, error)
	cancel(string) (string, error)
	noCommand(string, string, string, tgbotapi.FileID) (string, error)
	helpCommand(map[string]Command) (string, error)
	stateCommand(string, map[string]Command) (string, error)
	ValidateState(string, string) error
	ValidateAndInitUser(string) error
}

type Handler struct {
	collection    CollectionService
	trackerClient YandexTrackerService
}

func NewHandler(repo *infra.MongoRepository, trackerClient *api.Client) Receiver {
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
	b.WriteString("Ваши задачи:\n")
	for i, issue := range issues {
		b.WriteString(fmt.Sprintf("%v) %s: %s\n", i+1, issue.Key, issue.Link))
	}
	result := b.String()
	if len(issues) == 0 {
		result = "Вы ещё не создали ни одной задачи."
	}
	return result, nil
}

func (h Handler) createIssue(username, text string, commandMap map[string]Command) (string, error) {
	ctx := context.Background()

	if err := h.collection.UpdateStateUser(ctx, username, CREATING_STATE); err != nil {
		return "", err
	}

	if err := h.collection.ClearIssue(ctx, username); err != nil {
		return "", err
	}

	var b strings.Builder

	b.WriteString("Теперь вы можете создавать задачу.\n\n")
	b.WriteString("Используйте тэги для наполнения:\n")
	b.WriteString("\n--- Наполнение задачи ---\n\n")
	issue, err := h.collection.GetIssue(ctx, username)
	if err != nil {
		return "", err
	}
	b.WriteString(GetLocalizedIssueFilling(issue))
	b.WriteString("\n-------------------------\n\n")
	b.WriteString(fmt.Sprintf("Для проверки состояния, используйте:\n%s\n\n", commandMap[STATE_COMMAND].GetDescription()))
	b.WriteString(fmt.Sprintf("Когда будете готовы, используйте:\n%s\n\n", commandMap[DONE_COMMAND].GetDescription()))
	b.WriteString(fmt.Sprintf("Если передумаете, используйте:\n%s\n\n", commandMap[CANCEL_COMMAND].GetDescription()))

	return b.String(), nil
}

func (h Handler) done(ctx context.Context, username string) (string, error) {
	dbCtx := context.Background()

	user, err := h.collection.GetUser(dbCtx, username)
	if err != nil {
		return "", err
	}

	err = user.ValidateRequest()
	if err != nil {
		return "", err
	}

	SendPreliminaryMessagesWithContext(ctx, "Я начал создание задачи, по готовности отпишусь.", nil)

	issueData, err := h.createTrackerIssue(dbCtx, user)
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

func (h Handler) noCommand(username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	if tag == "" {
		return "", ErrEmptyTag
	}
	var err error
	ctx := context.Background()

	switch tag {
	case ISSUE_SUMMARY_TAG:
		err = h.collection.UpdateSummaryIssue(ctx, username, text)
	case ISSUE_DESCRIPTION_TAG:
		err = h.collection.UpdateDescriptionIssue(ctx, username, text)
	case ISSUE_ATTACHMENT_TAG:
		if fileID != "" {
			err = h.collection.AppendAttachmentIssue(ctx, username, string(fileID))
		} else {
			err = ErrEmptyAttachment
		}
	case ISSUE_DESCRIPTION_ATTACHMENT_TAG:
		if fileID != "" {
			err = h.collection.AppendDescriptionAttachmentIssue(ctx, username, string(fileID))
		} else {
			err = ErrEmptyAttachment
		}
	case ISSUE_TAGS_TAG:
		re := regexp.MustCompile(`[,\s.;]+`)
		replaced := re.ReplaceAllString(text, " ")
		rawTags := strings.Split(replaced, " ")
		tags := make([]string, 0, len(rawTags))

		for _, tag := range rawTags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
		err = h.collection.AppendTagIssue(ctx, username, tags)
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Успешно обновил %s", LocalizedTagsDescriptionMap[tag]), nil
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
		b.WriteString("\n-------------------------")
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

func (h Handler) createTrackerIssue(dbCtx context.Context, user *User) (*Issue, error) {

	err := h.uploadAttachments(user)
	if err != nil {
		return nil, err
	}

	err = h.uploadDescriptionAttachments(user)
	if err != nil {
		return nil, err
	}

	created, err := h.trackerClient.CreateIssue(user.Issue)
	if err != nil {
		return nil, err
	}

	issueData := Issue{
		Key:  created.Key,
		Link: NewIssueLink(created.Key),
	}

	err = h.collection.AppendDataIssue(dbCtx, user.Username, &issueData)
	if err != nil {
		return nil, err
	}

	err = h.collection.ClearIssue(dbCtx, user.Username)
	if err != nil {
		return nil, err
	}

	err = h.collection.UpdateStateUser(dbCtx, user.Username, DONE_STATE)
	if err != nil {
		return nil, err
	}

	return &issueData, nil
}

func (h Handler) uploadAttachments(user *User) error {
	if len(user.Issue.AttachmentIds) > 0 {
		var newAttachmentIDs []string
		for _, fileID := range user.Issue.AttachmentIds {
			fileConfig := tgbotapi.FileConfig{FileID: fileID}
			file, err := telegram.GetFileByConfig(&fileConfig)
			if err != nil {
				return err
			}
			fileURL := telegram.GetFileURLByPath(file.FilePath)
			resp, err := http.Get(fileURL)
			if err != nil {
				return err
			}
			res, err := h.trackerClient.UploadTemporaryAttachment(&resty.MultipartField{
				Reader:   resp.Body,
				FileName: strings.Split(file.FilePath, "/")[1],
			})
			if err != nil {
				return err
			}
			newAttachmentIDs = append(newAttachmentIDs, res.ID)
		}
		user.Issue.AttachmentIds = newAttachmentIDs
	}
	return nil
}

func (h Handler) uploadDescriptionAttachments(user *User) error {
	if len(user.Issue.DescriptionAttachmentIds) > 0 {
		var newDescriptionAttachmentIDs []string
		for _, fileID := range user.Issue.DescriptionAttachmentIds {
			fileConfig := tgbotapi.FileConfig{FileID: fileID}
			file, err := telegram.GetFileByConfig(&fileConfig)
			if err != nil {
				return err
			}
			fileURL := telegram.GetFileURLByPath(file.FilePath)
			resp, err := http.Get(fileURL)
			if err != nil {
				return err
			}
			filename := strings.Split(file.FilePath, "/")[1]
			res, err := h.trackerClient.UploadTemporaryAttachment(&resty.MultipartField{
				Reader:   resp.Body,
				FileName: filename,
			})
			if err != nil {
				return err
			}
			newDescriptionAttachmentIDs = append(newDescriptionAttachmentIDs, res.ID)
			if strings.HasSuffix(res.Name, ".svg") {
				user.Issue.Description += fmt.Sprintf("\n\n[%s](/ajax/v2/attachments/%s)", res.Name, res.ID)
			} else {
				user.Issue.Description += fmt.Sprintf("\n\n![%s](/ajax/v2/attachments/%s?inline=true =250x250)", res.Name, res.ID)
			}
		}
		user.Issue.DescriptionAttachmentIds = newDescriptionAttachmentIDs
	}
	return nil
}
