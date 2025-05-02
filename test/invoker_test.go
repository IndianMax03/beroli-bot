package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/IndianMax03/beroli-bot/internal/domain"
	mocks "github.com/IndianMax03/beroli-bot/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type InputCommandDataWantExecutedCommand struct {
	Name          string
	InputContext  context.Context
	InputUsername string
	InputText     string
	InputFileID   tgbotapi.FileID
	WantResult    string
	WantError     error
}

type InputStringWantParsedCommand struct {
	Name        string
	Input       string
	WantCmdName string
	WantText    string
	WantTag     string
	WantError   error
}

func TestInvokerExecuteCommandPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiver := mocks.NewMockReceiver(ctrl)
	nilCommandStub := domain.NewNilCommand(mockReceiver)
	commandMap := map[string]domain.Command{
		domain.NIL_COMMAND: nilCommandStub,
	}

	mockReceiver.
		EXPECT().
		ValidateAndInitUser(usernameStub).
		Return(nil)

	mockReceiver.
		EXPECT().
		ValidateState(usernameStub, nilCommandStub.GetName()).
		Return(nil)

	cases := []InputCommandDataWantExecutedCommand{
		{
			Name:          "single existing command",
			InputContext:  context.TODO(),
			InputUsername: usernameStub,
			InputText:     textStub,
			InputFileID:   fileIDStub,
			WantResult:    textStub,
		},
	}

	invoker := &domain.Invoker{
		Commands: commandMap,
		Receiver: mockReceiver,
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mockReceiver.
				EXPECT().
				NoCommand(c.InputUsername, gomock.Any(), gomock.Any(), c.InputFileID).
				Return(textStub, nil)
			_, actErr := invoker.ExecuteCommand(c.InputContext, c.InputUsername, c.InputText, c.InputFileID)
			assert.NoError(t, actErr, "An error occured in positive case")
		})
	}
}

func TestInvokerExecuteCommandNegativeErrValidateState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiver := mocks.NewMockReceiver(ctrl)
	nilCommandStub := domain.NewNilCommand(mockReceiver)
	commandMap := map[string]domain.Command{
		domain.NIL_COMMAND: nilCommandStub,
	}

	mockReceiver.
		EXPECT().
		ValidateAndInitUser(usernameStub).
		Return(nil)

	mockReceiver.
		EXPECT().
		ValidateState(usernameStub, nilCommandStub.GetName()).
		Return(errStub)

	testCase := InputCommandDataWantExecutedCommand{
		Name:          "single existing command with error during ValidateState",
		InputContext:  context.TODO(),
		InputUsername: usernameStub,
		InputText:     textStub,
		InputFileID:   fileIDStub,
		WantResult:    textStub,
	}

	invoker := &domain.Invoker{
		Commands: commandMap,
		Receiver: mockReceiver,
	}

	t.Run(testCase.Name, func(t *testing.T) {
		_, actErr := invoker.ExecuteCommand(testCase.InputContext, testCase.InputUsername, testCase.InputText, testCase.InputFileID)
		assert.ErrorIs(t, actErr, errStub)
	})
}

func TestInvokerExecuteCommandNegativeErrValidateAndInitUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiver := mocks.NewMockReceiver(ctrl)
	nilCommandStub := domain.NewNilCommand(mockReceiver)
	commandMap := map[string]domain.Command{
		domain.NIL_COMMAND: nilCommandStub,
	}

	mockReceiver.
		EXPECT().
		ValidateAndInitUser(usernameStub).
		Return(errStub)

	testCase := InputCommandDataWantExecutedCommand{
		Name:          "single existing command with error during ValidateAndInitUser",
		InputContext:  context.TODO(),
		InputUsername: usernameStub,
		InputText:     textStub,
		InputFileID:   fileIDStub,
		WantResult:    textStub,
	}

	invoker := &domain.Invoker{
		Commands: commandMap,
		Receiver: mockReceiver,
	}

	t.Run(testCase.Name, func(t *testing.T) {
		_, actErr := invoker.ExecuteCommand(testCase.InputContext, testCase.InputUsername, testCase.InputText, testCase.InputFileID)
		assert.ErrorIs(t, actErr, errStub)
	})
}

func TestInvokerParseCommandPositive(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiver := mocks.NewMockReceiver(ctrl)

	cases := []InputStringWantParsedCommand{
		{
			Name:        "single existing command",
			Input:       domain.CREATE_ISSUE_COMMAND,
			WantCmdName: domain.CREATE_ISSUE_COMMAND,
			WantText:    "",
			WantTag:     "",
		},
		{
			Name:        "tag with text",
			Input:       fmt.Sprintf("%s %s", domain.ISSUE_ATTACHMENT_TAG, textStub),
			WantCmdName: domain.NIL_COMMAND,
			WantText:    textStub,
			WantTag:     domain.ISSUE_ATTACHMENT_TAG,
		},
		{
			Name:        "tag with newline text",
			Input:       fmt.Sprintf("%s\n%s", domain.ISSUE_ATTACHMENT_TAG, textStub),
			WantCmdName: domain.NIL_COMMAND,
			WantText:    textStub,
			WantTag:     domain.ISSUE_ATTACHMENT_TAG,
		},
		{
			Name:        "tag without text",
			Input:       domain.ISSUE_ATTACHMENT_TAG,
			WantCmdName: domain.NIL_COMMAND,
			WantText:    "",
			WantTag:     domain.ISSUE_ATTACHMENT_TAG,
		},
		{
			Name:        "single text",
			Input:       textStub,
			WantCmdName: domain.NIL_COMMAND,
			WantText:    textStub,
			WantTag:     "",
		},
		{
			Name:        "single text with spaces and newline",
			Input:       fmt.Sprintf("%s %s %s.\n%s %s", textStub, textStub, textStub, textStub, textStub),
			WantCmdName: domain.NIL_COMMAND,
			WantText:    fmt.Sprintf("%s %s %s.\n%s %s", textStub, textStub, textStub, textStub, textStub),
			WantTag:     "",
		},
		{
			Name:        "empty text",
			Input:       "",
			WantCmdName: domain.NIL_COMMAND,
			WantText:    "",
			WantTag:     "",
		},
	}

	commandMap := map[string]domain.Command{
		domain.NIL_COMMAND: domain.NewNilCommand(mockReceiver),
	}
	commandMap[domain.CREATE_ISSUE_COMMAND] = domain.NewCreateIssueCommand(mockReceiver, commandMap)

	invoker := &domain.Invoker{
		Commands: commandMap,
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			actCmd, actText, actTag, actErr := invoker.ParseCommand(c.Input)
			assert.NoError(t, actErr, "An error occured in positive case")
			assert.Equal(t, c.WantCmdName, actCmd.GetName(), "Wrong command name received")
			assert.Equal(t, c.WantText, actText, "Wrong text received")
			assert.Equal(t, c.WantTag, actTag, "Wrong tag received")
		})
	}
}

func TestInvokerParseCommandNegative(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiver := mocks.NewMockReceiver(ctrl)

	cases := []InputStringWantParsedCommand{
		{
			Name:      "multiply spaced tags",
			Input:     fmt.Sprintf("%s %s", domain.ISSUE_ATTACHMENT_TAG, domain.ISSUE_SUMMARY_TAG),
			WantError: domain.ErrMultiplyTags,
		},
		{
			Name:      "multiply newlined tags",
			Input:     fmt.Sprintf("%s\n%s", domain.ISSUE_ATTACHMENT_TAG, domain.ISSUE_SUMMARY_TAG),
			WantError: domain.ErrMultiplyTags,
		},
		{
			Name:      "unknown command",
			Input:     fmt.Sprintf("/%s", RandStringWithSymbolsAndEmojis(5)),
			WantError: domain.ErrUnknownCommand,
		},
		{
			Name:      "unknown tag",
			Input:     fmt.Sprintf("#%s", RandStringWithSymbolsAndEmojis(5)),
			WantError: domain.ErrUnknownTag,
		},
	}

	commandMap := map[string]domain.Command{
		domain.NIL_COMMAND: domain.NewNilCommand(mockReceiver),
	}
	commandMap[domain.CREATE_ISSUE_COMMAND] = domain.NewCreateIssueCommand(mockReceiver, commandMap)

	invoker := &domain.Invoker{
		Commands: commandMap,
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, _, _, actErr := invoker.ParseCommand(c.Input)
			assert.ErrorIs(t, actErr, c.WantError)
		})
	}
}
