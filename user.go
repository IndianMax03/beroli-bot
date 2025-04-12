package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
)

const (
	USERNAME_FIELD      = "username"
	STATE_FIELD         = "state"
	ISSUE_FIELD         = "issue"
	ISSUES_FIELD        = "issues"
	TRACKER_DEFAULT_URL = "https://tracker.yandex.ru/"

	ISSUE_SUMMARY_TAG                = "#sum"
	ISSUE_DESCRIPTION_TAG            = "#des"
	ISSUE_ATTACHMENT_TAG             = "#att"
	ISSUE_DESCRIPTION_ATTACHMENT_TAG = "#datt"
	ISSUE_TAGS_TAG                   = "#tags"
)

var (
	ErrMandatoryIsRequired = errors.New("mandatory fields is required")
)

var (
	LocalizedTagsDescriptionMap = map[string]string{
		ISSUE_SUMMARY_TAG:                "название",
		ISSUE_DESCRIPTION_TAG:            "описание",
		ISSUE_ATTACHMENT_TAG:             "вложения",
		ISSUE_DESCRIPTION_ATTACHMENT_TAG: "вложения к описанию",
		ISSUE_TAGS_TAG:                   "тэги",
	}
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

func NewIssueLink(key string) string {
	return TRACKER_DEFAULT_URL + key
}

func NewDefaultIssue() *model.IssueCreateRequest {
	return &model.IssueCreateRequest{
		Queue: model.Queue{
			Key: TRACKER_QUEUE,
		},
		Type:                     "bug",
		Priority:                 model.CriticalPriority,
		AttachmentIds:            []string{},
		DescriptionAttachmentIds: []string{},
		Tags:                     []string{},
	}
}

func (u *User) validateRequest() error {
	if u.Issue.Summary == "" || u.Issue.Queue.Key == "" {
		return ErrMandatoryIsRequired
	}
	return nil
}

func GetLocalizedIssueFilling(issue *model.IssueCreateRequest) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s - %s\n", ISSUE_SUMMARY_TAG, LocalizedTagsDescriptionMap[ISSUE_SUMMARY_TAG]))
	if issue.Summary != "" {
		b.WriteString(fmt.Sprintf("%s\n", CutString(issue.Summary)))
	} else {
		b.WriteString("‼️ Обязательное поле не заполнено ‼️\n")
	}

	b.WriteString(fmt.Sprintf("%s - %s\n", ISSUE_DESCRIPTION_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_TAG]))
	if issue.Description != "" {
		b.WriteString(fmt.Sprintf("%s\n", CutString(issue.Description)))
	} else {
		b.WriteString("Отсутствует\n")
	}

	b.WriteString(fmt.Sprintf("%s - %s\n", ISSUE_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_ATTACHMENT_TAG]))
	if len(issue.AttachmentIds) > 0 {
		b.WriteString(fmt.Sprintf("%s\n", CutArrayOfString(issue.AttachmentIds)))
	} else {
		b.WriteString("Отсутствуют\n")
	}

	b.WriteString(fmt.Sprintf("%s - %s\n", ISSUE_DESCRIPTION_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_ATTACHMENT_TAG]))
	if len(issue.DescriptionAttachmentIds) > 0 {
		b.WriteString(fmt.Sprintf("%s\n", CutArrayOfString(issue.DescriptionAttachmentIds)))
	} else {
		b.WriteString("Отсутствуют\n")
	}

	b.WriteString(fmt.Sprintf("%s - %s\n", ISSUE_TAGS_TAG, LocalizedTagsDescriptionMap[ISSUE_TAGS_TAG]))
	if len(issue.Tags) > 0 {
		b.WriteString(fmt.Sprintf("%s\n", CutArrayOfString(issue.Tags)))
	} else {
		b.WriteString("Отсутствуют\n")
	}

	return b.String()
}
