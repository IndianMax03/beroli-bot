package test

import (
	"testing"

	"github.com/IndianMax03/beroli-bot/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

type InputMessageWantParsedMessage struct {
	Name    string
	Input   *tgbotapi.Message
	WantErr error
}

func TestParserParseMessagePositive(t *testing.T) {
	cases := []InputMessageWantParsedMessage{
		{
			Name: "text from user",
			Input: &tgbotapi.Message{
				From: &tgbotapi.User{
					UserName: usernameStub,
				},
				Caption:  "",
				Text:     textStub,
				Document: nil,
				Photo:    nil,
			},
		},
		{
			Name: "valid size document from user",
			Input: &tgbotapi.Message{
				From: &tgbotapi.User{
					UserName: usernameStub,
				},
				Caption: textStub,
				Text:    "",
				Document: &tgbotapi.Document{
					FileSize: util.MAX_ATTACHMENT_SIZE_BYTE,
					FileID:   textStub,
				},
				Photo: nil,
			},
		},
		{
			Name: "valid size photo from user",
			Input: &tgbotapi.Message{
				From: &tgbotapi.User{
					UserName: usernameStub,
				},
				Caption:  textStub,
				Text:     "",
				Document: nil,
				Photo: []tgbotapi.PhotoSize{
					{
						FileSize: util.MAX_ATTACHMENT_SIZE_BYTE,
						FileID:   textStub,
					},
				},
			},
		},
		{
			Name: "valid size photo from user without caption",
			Input: &tgbotapi.Message{
				From: &tgbotapi.User{
					UserName: usernameStub,
				},
				Caption:  "",
				Text:     "",
				Document: nil,
				Photo: []tgbotapi.PhotoSize{
					{
						FileSize: util.MAX_ATTACHMENT_SIZE_BYTE,
						FileID:   textStub,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actUsername, actText, actFileID, actErr := util.ParseMessage(c.Input)
			assert.NoError(t, actErr, "An error occured in positive case")
			assert.Equal(t, c.Input.From.UserName, actUsername)
			if c.Input.Document != nil || c.Input.Photo != nil {
				assert.Equal(t, c.Input.Caption, actText)
				assert.NotEmpty(t, actFileID)
			} else {
				assert.Equal(t, c.Input.Text, actText)
				assert.Empty(t, actFileID)
			}
		})
	}
}

func TestParserParseMessageNegative(t *testing.T) {
	cases := []InputMessageWantParsedMessage{
		{
			Name: "invalid size document from user",
			Input: &tgbotapi.Message{
				From: &tgbotapi.User{
					UserName: usernameStub,
				},
				Caption:  textStub,
				Text:     "",
				Document: nil,
				Photo: []tgbotapi.PhotoSize{
					{
						FileSize: util.MAX_ATTACHMENT_SIZE_BYTE + 1,
						FileID:   textStub,
					},
				},
			},
			WantErr: util.ErrAttachmentSize,
		},
		{
			Name: "invalid size photo from user",
			Input: &tgbotapi.Message{
				From: &tgbotapi.User{
					UserName: usernameStub,
				},
				Caption:  textStub,
				Text:     "",
				Document: nil,
				Photo: []tgbotapi.PhotoSize{
					{
						FileSize: util.MAX_ATTACHMENT_SIZE_BYTE + 1,
						FileID:   textStub,
					},
				},
			},
			WantErr: util.ErrAttachmentSize,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, _, _, actErr := util.ParseMessage(c.Input)
			assert.ErrorIs(t, actErr, c.WantErr)
		})
	}
}
