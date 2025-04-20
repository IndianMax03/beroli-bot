package main

import (
	"fmt"
	"testing"

	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	"github.com/stretchr/testify/assert"
)

type InputUserWantError struct {
	InputUser *User
	Error     error
}

type InputIssueCreateRequestWantStringContains struct {
	Request      *model.IssueCreateRequest
	WantContains string
}

func generateCasesNewIssueLink() []InputStringWantString {
	strLen1 := RandStringWithSymbolsAndEmojis(1)
	strLen5 := RandStringWithSymbolsAndEmojis(5)
	strLen20 := RandStringWithSymbolsAndEmojis(20)
	cases := []InputStringWantString{
		{
			Input: "",
			Want:  TRACKER_DEFAULT_URL,
		},
		{
			Input: strLen1,
			Want:  TRACKER_DEFAULT_URL + strLen1,
		},
		{
			Input: strLen5,
			Want:  TRACKER_DEFAULT_URL + strLen5,
		},
		{
			Input: strLen20,
			Want:  TRACKER_DEFAULT_URL + strLen20,
		},
	}
	return cases
}

func generateCases_validateRequest() []InputUserWantError {
	TRACKER_QUEUE = ""

	issueEmptySummaryEmptyQueueKey := NewDefaultIssue()
	issueEmptySummaryEmptyQueueKey.Summary = ""
	userWithIssueEmptySummaryEmptyQueueKey := User{
		Issue: issueEmptySummaryEmptyQueueKey,
	}

	issueNonEmptySummaryEmptyQueueKey := NewDefaultIssue()
	issueNonEmptySummaryEmptyQueueKey.Summary = RandStringWithSymbolsAndEmojis(20)
	userWithIssueNonEmptySummaryEmptyQueueKey := User{
		Issue: issueNonEmptySummaryEmptyQueueKey,
	}

	TRACKER_QUEUE = RandStringWithSymbolsAndEmojis(5)

	issueEmptySummaryNonEmptyQueueKey := NewDefaultIssue()
	issueEmptySummaryNonEmptyQueueKey.Summary = ""
	userWithIssueEmptySummaryNonEmptyQueueKey := User{
		Issue: issueEmptySummaryNonEmptyQueueKey,
	}

	issue1LenNonEmptySummaryNonEmptyQueueKey := NewDefaultIssue()
	issue1LenNonEmptySummaryNonEmptyQueueKey.Summary = RandStringWithSymbolsAndEmojis(1)
	userWithIssue1LenNonEmptySummaryNonEmptyQueueKey := User{
		Issue: issue1LenNonEmptySummaryNonEmptyQueueKey,
	}

	issue20LenNonEmptySummaryNonEmptyQueueKey := NewDefaultIssue()
	issue20LenNonEmptySummaryNonEmptyQueueKey.Summary = RandStringWithSymbolsAndEmojis(20)
	userWithIssue20LenNonEmptySummaryNonEmptyQueueKey := User{
		Issue: issue20LenNonEmptySummaryNonEmptyQueueKey,
	}

	cases := []InputUserWantError{
		{
			InputUser: &userWithIssueEmptySummaryEmptyQueueKey,
			Error:     ErrMandatoryIsRequired,
		},
		{
			InputUser: &userWithIssueNonEmptySummaryEmptyQueueKey,
			Error:     ErrMandatoryIsRequired,
		},
		{
			InputUser: &userWithIssueEmptySummaryNonEmptyQueueKey,
			Error:     ErrMandatoryIsRequired,
		},
		{
			InputUser: &userWithIssue1LenNonEmptySummaryNonEmptyQueueKey,
			Error:     nil,
		},
		{
			InputUser: &userWithIssue20LenNonEmptySummaryNonEmptyQueueKey,
			Error:     nil,
		},
	}
	return cases
}

func TestNewIssueLink(t *testing.T) {
	tests := generateCasesNewIssueLink()
	for _, test := range tests {
		name := fmt.Sprintf("CASE:'%s'->'%s'", test.Input, test.Want)
		t.Run(name, func(t *testing.T) {
			got := NewIssueLink(test.Input)
			if got != test.Want {
				t.Errorf("Got: '%s', Want '%s'", got, test.Want)
			}
		})
	}
}

func TestNewDefaultIssue(t *testing.T) {
	assert := assert.New(t)
	TRACKER_QUEUE = RandStringWithSymbolsAndEmojis(5)
	emptyArrayOfString := []string{}

	issue := NewDefaultIssue()

	assert.Equal(TRACKER_QUEUE, issue.Queue.Key, fmt.Sprintf("Queue key must be equal to global value TRACKER_QUEUE: '%s'", TRACKER_QUEUE))
	assert.Equal(DEFAULT_ISSUE_TYPE, issue.Type, fmt.Sprintf("Type must be equal to the default value: '%s'", DEFAULT_ISSUE_TYPE))
	assert.Equal(DEFAULT_ISSUE_PRIORITY, issue.Priority, fmt.Sprintf("Priority must be equal to the default value: '%v'", DEFAULT_ISSUE_PRIORITY))
	assert.Equal(emptyArrayOfString, issue.AttachmentIds, "AttachmentIds must be initialized to empty string array")
	assert.Equal(emptyArrayOfString, issue.DescriptionAttachmentIds, "DescriptionAttachmentIds must be initialized to empty string array")
	assert.Equal(emptyArrayOfString, issue.Tags, "Tags must be initialized to empty string array")
}

func Test_validateRequest(t *testing.T) {
	assert := assert.New(t)
	tests := generateCases_validateRequest()

	for _, test := range tests {
		name := fmt.Sprintf("CASE:'(N:%s, QK:%s)' -> Error: '%s'", test.InputUser.Issue.Summary, test.InputUser.Issue.Queue.Key, test.Error)
		t.Run(name, func(t *testing.T) {
			err := test.InputUser.validateRequest()
			assert.Equal(test.Error, err)
		})
	}

}

func TestGetLocalizedIssueFillingSummary(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := NewDefaultIssue()
	randStringLenLessThanMaxStringLen := RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1)
	request.Summary = randStringLenLessThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_SUMMARY_TAG, LocalizedTagsDescriptionMap[ISSUE_SUMMARY_TAG], CutString(randStringLenLessThanMaxStringLen)),
	})

	randStringLenGreaterThanMaxStringLen := RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT + 1)
	request = NewDefaultIssue()
	request.Summary = randStringLenGreaterThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_SUMMARY_TAG, LocalizedTagsDescriptionMap[ISSUE_SUMMARY_TAG], CutString(randStringLenGreaterThanMaxStringLen)),
	})

	emptyString := ""
	request = NewDefaultIssue()
	request.Summary = emptyString
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n‼️ Обязательное поле не заполнено ‼️", ISSUE_SUMMARY_TAG, LocalizedTagsDescriptionMap[ISSUE_SUMMARY_TAG]),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Summary='%s'", test.Request.Summary)
		t.Run(name, func(t *testing.T) {
			res := GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingDescription(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := NewDefaultIssue()
	randStringLenLessThanMaxStringLen := RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1)
	request.Description = randStringLenLessThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_DESCRIPTION_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_TAG], CutString(randStringLenLessThanMaxStringLen)),
	})

	randStringLenGreaterThanMaxStringLen := RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT + 1)
	request = NewDefaultIssue()
	request.Description = randStringLenGreaterThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_DESCRIPTION_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_TAG], CutString(randStringLenGreaterThanMaxStringLen)),
	})

	emptyString := ""
	request = NewDefaultIssue()
	request.Description = emptyString
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\nОтсутствует", ISSUE_DESCRIPTION_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_TAG]),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Description='%s'", test.WantContains)
		t.Run(name, func(t *testing.T) {
			res := GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingAttachmentIds(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := NewDefaultIssue()
	emptyAttachmentIds := []string{}
	request.AttachmentIds = emptyAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_ATTACHMENT_TAG], "Отсутствуют"),
	})

	request = NewDefaultIssue()
	len1AttachmentIds := []string{RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1)}
	request.AttachmentIds = len1AttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к задаче", ISSUE_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_ATTACHMENT_TAG], len(len1AttachmentIds)),
	})

	request = NewDefaultIssue()
	len5AttachmentIds := []string{
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
	}
	assert.LessOrEqual(len(len5AttachmentIds), MAX_ARRAY_LENGHT, "Generated invalid case!")
	request.AttachmentIds = len5AttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к задаче", ISSUE_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_ATTACHMENT_TAG], len(len5AttachmentIds)),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Len AttachmentIds='%v'", len(test.Request.AttachmentIds))
		t.Run(name, func(t *testing.T) {
			res := GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingDescriptionAttachmentIds(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := NewDefaultIssue()
	emptyDescriptionAttachmentIds := []string{}
	request.DescriptionAttachmentIds = emptyDescriptionAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_DESCRIPTION_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_ATTACHMENT_TAG], "Отсутствуют"),
	})

	request = NewDefaultIssue()
	len1DescriptionAttachmentIds := []string{RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1)}
	request.DescriptionAttachmentIds = len1DescriptionAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к описанию", ISSUE_DESCRIPTION_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_ATTACHMENT_TAG], len(len1DescriptionAttachmentIds)),
	})

	request = NewDefaultIssue()
	len5DescriptionAttachmentIds := []string{
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
	}
	assert.LessOrEqual(len(len5DescriptionAttachmentIds), MAX_ARRAY_LENGHT, "Generated invalid case!")
	request.DescriptionAttachmentIds = len5DescriptionAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к описанию", ISSUE_DESCRIPTION_ATTACHMENT_TAG, LocalizedTagsDescriptionMap[ISSUE_DESCRIPTION_ATTACHMENT_TAG], len(len5DescriptionAttachmentIds)),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Len DescriptionAttachmentIds='%v'", len(test.Request.DescriptionAttachmentIds))
		t.Run(name, func(t *testing.T) {
			res := GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingTags(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := NewDefaultIssue()
	emptyTags := []string{}
	request.Tags = emptyTags
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_TAGS_TAG, LocalizedTagsDescriptionMap[ISSUE_TAGS_TAG], "Отсутствуют"),
	})

	request = NewDefaultIssue()
	len1Tags := []string{RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1)}
	request.Tags = len1Tags
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_TAGS_TAG, LocalizedTagsDescriptionMap[ISSUE_TAGS_TAG], CutArrayOfString(len1Tags)),
	})

	request = NewDefaultIssue()
	len5Tags := []string{
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(MAX_STRING_LENGHT - 1),
	}
	assert.LessOrEqual(len(len5Tags), MAX_ARRAY_LENGHT, "Generated invalid case!")
	request.Tags = len5Tags
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", ISSUE_TAGS_TAG, LocalizedTagsDescriptionMap[ISSUE_TAGS_TAG], CutArrayOfString(len5Tags)),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Len Tags='%v'", len(test.Request.Tags))
		t.Run(name, func(t *testing.T) {
			res := GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}
