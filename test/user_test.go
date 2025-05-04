package test

import (
	"fmt"
	"testing"

	"github.com/IndianMax03/beroli-bot/internal/domain"
	"github.com/IndianMax03/beroli-bot/internal/global"
	"github.com/IndianMax03/beroli-bot/internal/util"
	"github.com/IndianMax03/yandex-tracker-go-client/v3/model"
	"github.com/stretchr/testify/assert"
)

type InputUserWantError struct {
	InputUser *domain.User
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
			Want:  domain.TRACKER_DEFAULT_URL,
		},
		{
			Input: strLen1,
			Want:  domain.TRACKER_DEFAULT_URL + strLen1,
		},
		{
			Input: strLen5,
			Want:  domain.TRACKER_DEFAULT_URL + strLen5,
		},
		{
			Input: strLen20,
			Want:  domain.TRACKER_DEFAULT_URL + strLen20,
		},
	}
	return cases
}

func generateCases_validateRequest() []InputUserWantError {
	global.TRACKER_QUEUE = ""

	issueEmptySummaryEmptyQueueKey := domain.NewDefaultIssue()
	issueEmptySummaryEmptyQueueKey.Summary = ""
	userWithIssueEmptySummaryEmptyQueueKey := domain.User{
		Issue: issueEmptySummaryEmptyQueueKey,
	}

	issueNonEmptySummaryEmptyQueueKey := domain.NewDefaultIssue()
	issueNonEmptySummaryEmptyQueueKey.Summary = RandStringWithSymbolsAndEmojis(20)
	userWithIssueNonEmptySummaryEmptyQueueKey := domain.User{
		Issue: issueNonEmptySummaryEmptyQueueKey,
	}

	global.TRACKER_QUEUE = RandStringWithSymbolsAndEmojis(5)

	issueEmptySummaryNonEmptyQueueKey := domain.NewDefaultIssue()
	issueEmptySummaryNonEmptyQueueKey.Summary = ""
	userWithIssueEmptySummaryNonEmptyQueueKey := domain.User{
		Issue: issueEmptySummaryNonEmptyQueueKey,
	}

	issue1LenNonEmptySummaryNonEmptyQueueKey := domain.NewDefaultIssue()
	issue1LenNonEmptySummaryNonEmptyQueueKey.Summary = RandStringWithSymbolsAndEmojis(1)
	userWithIssue1LenNonEmptySummaryNonEmptyQueueKey := domain.User{
		Issue: issue1LenNonEmptySummaryNonEmptyQueueKey,
	}

	issue20LenNonEmptySummaryNonEmptyQueueKey := domain.NewDefaultIssue()
	issue20LenNonEmptySummaryNonEmptyQueueKey.Summary = RandStringWithSymbolsAndEmojis(20)
	userWithIssue20LenNonEmptySummaryNonEmptyQueueKey := domain.User{
		Issue: issue20LenNonEmptySummaryNonEmptyQueueKey,
	}

	cases := []InputUserWantError{
		{
			InputUser: &userWithIssueEmptySummaryEmptyQueueKey,
			Error:     domain.ErrMandatoryIsRequired,
		},
		{
			InputUser: &userWithIssueNonEmptySummaryEmptyQueueKey,
			Error:     domain.ErrMandatoryIsRequired,
		},
		{
			InputUser: &userWithIssueEmptySummaryNonEmptyQueueKey,
			Error:     domain.ErrMandatoryIsRequired,
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
			got := domain.NewIssueLink(test.Input)
			if got != test.Want {
				t.Errorf("Got: '%s', Want '%s'", got, test.Want)
			}
		})
	}
}

func TestNewDefaultIssue(t *testing.T) {
	assert := assert.New(t)
	global.TRACKER_QUEUE = RandStringWithSymbolsAndEmojis(5)
	emptyArrayOfString := []string{}

	issue := domain.NewDefaultIssue()

	assert.Equal(global.TRACKER_QUEUE, issue.Queue.Key, fmt.Sprintf("Queue key must be equal to global value TRACKER_QUEUE: '%s'", global.TRACKER_QUEUE))
	assert.Equal(domain.DEFAULT_ISSUE_TYPE, issue.Type, fmt.Sprintf("Type must be equal to the default value: '%s'", domain.DEFAULT_ISSUE_TYPE))
	assert.Equal(domain.DEFAULT_ISSUE_PRIORITY, issue.Priority, fmt.Sprintf("Priority must be equal to the default value: '%v'", domain.DEFAULT_ISSUE_PRIORITY))
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
			err := test.InputUser.ValidateRequest()
			assert.Equal(test.Error, err)
		})
	}

}

func TestGetLocalizedIssueFillingSummary(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := domain.NewDefaultIssue()
	randStringLenLessThanMaxStringLen := RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1)
	request.Summary = randStringLenLessThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_SUMMARY_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_SUMMARY_TAG], util.CutString(randStringLenLessThanMaxStringLen)),
	})

	randStringLenGreaterThanMaxStringLen := RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT + 1)
	request = domain.NewDefaultIssue()
	request.Summary = randStringLenGreaterThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_SUMMARY_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_SUMMARY_TAG], util.CutString(randStringLenGreaterThanMaxStringLen)),
	})

	emptyString := ""
	request = domain.NewDefaultIssue()
	request.Summary = emptyString
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n‼️ Обязательное поле не заполнено ‼️", domain.ISSUE_SUMMARY_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_SUMMARY_TAG]),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Summary='%s'", test.Request.Summary)
		t.Run(name, func(t *testing.T) {
			res := domain.GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingDescription(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := domain.NewDefaultIssue()
	randStringLenLessThanMaxStringLen := RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1)
	request.Description = randStringLenLessThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_DESCRIPTION_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_DESCRIPTION_TAG], util.CutString(randStringLenLessThanMaxStringLen)),
	})

	randStringLenGreaterThanMaxStringLen := RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT + 1)
	request = domain.NewDefaultIssue()
	request.Description = randStringLenGreaterThanMaxStringLen
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_DESCRIPTION_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_DESCRIPTION_TAG], util.CutString(randStringLenGreaterThanMaxStringLen)),
	})

	emptyString := ""
	request = domain.NewDefaultIssue()
	request.Description = emptyString
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\nОтсутствует", domain.ISSUE_DESCRIPTION_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_DESCRIPTION_TAG]),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Description='%s'", test.WantContains)
		t.Run(name, func(t *testing.T) {
			res := domain.GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingAttachmentIds(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := domain.NewDefaultIssue()
	emptyAttachmentIds := []string{}
	request.AttachmentIds = emptyAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_ATTACHMENT_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_ATTACHMENT_TAG], "Отсутствуют"),
	})

	request = domain.NewDefaultIssue()
	len1AttachmentIds := []string{RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1)}
	request.AttachmentIds = len1AttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к задаче", domain.ISSUE_ATTACHMENT_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_ATTACHMENT_TAG], len(len1AttachmentIds)),
	})

	request = domain.NewDefaultIssue()
	len5AttachmentIds := []string{
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
	}
	assert.LessOrEqual(len(len5AttachmentIds), util.MAX_ARRAY_LENGHT, "Generated invalid case!")
	request.AttachmentIds = len5AttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к задаче", domain.ISSUE_ATTACHMENT_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_ATTACHMENT_TAG], len(len5AttachmentIds)),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Len AttachmentIds='%v'", len(test.Request.AttachmentIds))
		t.Run(name, func(t *testing.T) {
			res := domain.GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingDescriptionAttachmentIds(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := domain.NewDefaultIssue()
	emptyDescriptionAttachmentIds := []string{}
	request.DescriptionAttachmentIds = emptyDescriptionAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG], "Отсутствуют"),
	})

	request = domain.NewDefaultIssue()
	len1DescriptionAttachmentIds := []string{RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1)}
	request.DescriptionAttachmentIds = len1DescriptionAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к описанию", domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG], len(len1DescriptionAttachmentIds)),
	})

	request = domain.NewDefaultIssue()
	len5DescriptionAttachmentIds := []string{
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
	}
	assert.LessOrEqual(len(len5DescriptionAttachmentIds), util.MAX_ARRAY_LENGHT, "Generated invalid case!")
	request.DescriptionAttachmentIds = len5DescriptionAttachmentIds
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%v вложений к описанию", domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_DESCRIPTION_ATTACHMENT_TAG], len(len5DescriptionAttachmentIds)),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Len DescriptionAttachmentIds='%v'", len(test.Request.DescriptionAttachmentIds))
		t.Run(name, func(t *testing.T) {
			res := domain.GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}

func TestGetLocalizedIssueFillingTags(t *testing.T) {
	assert := assert.New(t)
	var cases []InputIssueCreateRequestWantStringContains

	request := domain.NewDefaultIssue()
	emptyTags := []string{}
	request.Tags = emptyTags
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_TAGS_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_TAGS_TAG], "Отсутствуют"),
	})

	request = domain.NewDefaultIssue()
	len1Tags := []string{RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1)}
	request.Tags = len1Tags
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_TAGS_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_TAGS_TAG], util.CutArrayOfString(len1Tags)),
	})

	request = domain.NewDefaultIssue()
	len5Tags := []string{
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
		RandStringWithSymbolsAndEmojis(util.MAX_STRING_LENGHT - 1),
	}
	assert.LessOrEqual(len(len5Tags), util.MAX_ARRAY_LENGHT, "Generated invalid case!")
	request.Tags = len5Tags
	cases = append(cases, InputIssueCreateRequestWantStringContains{
		Request:      request,
		WantContains: fmt.Sprintf("%s - %s\n%s", domain.ISSUE_TAGS_TAG, domain.LocalizedTagsDescriptionMap[domain.ISSUE_TAGS_TAG], util.CutArrayOfString(len5Tags)),
	})

	for _, test := range cases {
		name := fmt.Sprintf("CASE: Len Tags='%v'", len(test.Request.Tags))
		t.Run(name, func(t *testing.T) {
			res := domain.GetLocalizedIssueFilling(test.Request)
			assert.Contains(res, test.WantContains)
		})
	}
}
