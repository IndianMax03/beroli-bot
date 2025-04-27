package domain

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MY_ISSUES_COMMAND             = "/my_issues"
	MY_ISSUES_COMMAND_DESCRIPTION = "получить список созданных задач"

	CREATE_ISSUE_COMMAND             = "/issue"
	CREATE_ISSUE_COMMAND_DESCRIPTION = "начать создание новой задачи"

	DONE_COMMAND             = "/done"
	DONE_COMMAND_DESCRIPTION = "завершить создание задачи"

	CANCEL_COMMAND             = "/cancel"
	CANCEL_COMMAND_DESCRIPTION = "отменить создание задачи"

	NIL_COMMAND             = ""
	NIL_COMMAND_DESCRIPTION = ""

	HELP_COMMAND             = "/help"
	HELP_COMMAND_DESCRIPTION = "получить справку по командам"

	STATE_COMMAND             = "/state"
	STATE_COMMAND_DESCRIPTION = "узнать текущий контекст"
)

func NewMyIssuesCommand(r Receiver) *MyIssuesCommand {
	return &MyIssuesCommand{
		concreteCommand: concreteCommand{
			commandName:        MY_ISSUES_COMMAND,
			commandDescription: MY_ISSUES_COMMAND_DESCRIPTION,
		},
		receiver: r,
	}
}

func NewCreateIssueCommand(r Receiver, commandMap map[string]Command) *CreateIssueCommand {
	return &CreateIssueCommand{
		concreteCommand: concreteCommand{
			commandName:        CREATE_ISSUE_COMMAND,
			commandDescription: CREATE_ISSUE_COMMAND_DESCRIPTION,
		},
		receiver:   r,
		commandMap: commandMap,
	}
}

func NewDoneCommand(r Receiver) *DoneCommand {
	return &DoneCommand{
		concreteCommand: concreteCommand{
			commandName:        DONE_COMMAND,
			commandDescription: DONE_COMMAND_DESCRIPTION,
		},
		receiver: r,
	}
}

func NewCancelCommand(r Receiver) *CancelCommand {
	return &CancelCommand{
		concreteCommand: concreteCommand{
			commandName:        CANCEL_COMMAND,
			commandDescription: CANCEL_COMMAND_DESCRIPTION,
		},
		receiver: r,
	}
}

func NewNilCommand(r Receiver) *NilCommand {
	return &NilCommand{
		concreteCommand: concreteCommand{
			commandName:        NIL_COMMAND,
			commandDescription: NIL_COMMAND_DESCRIPTION,
		},
		receiver: r,
	}
}

func NewHelpCommand(r Receiver, commandMap map[string]Command) *HelpCommand {
	return &HelpCommand{
		concreteCommand: concreteCommand{
			commandName:        HELP_COMMAND,
			commandDescription: HELP_COMMAND_DESCRIPTION,
		},
		receiver:   r,
		commandMap: commandMap,
	}
}

func NewStateCommand(r Receiver, commandMap map[string]Command) *StateCommand {
	return &StateCommand{
		concreteCommand: concreteCommand{
			commandName:        STATE_COMMAND,
			commandDescription: STATE_COMMAND_DESCRIPTION,
		},
		receiver:   r,
		commandMap: commandMap,
	}
}

type Command interface {
	execute(context.Context, string, string, string, tgbotapi.FileID) (string, error)
	GetName() string
	GetDescription() string
}

type MyIssuesCommand struct {
	concreteCommand
	receiver Receiver
}

type CreateIssueCommand struct {
	concreteCommand
	receiver   Receiver
	commandMap map[string]Command
}

type DoneCommand struct {
	concreteCommand
	receiver Receiver
}

type CancelCommand struct {
	concreteCommand
	receiver Receiver
}

type NilCommand struct {
	concreteCommand
	receiver Receiver
}

type HelpCommand struct {
	concreteCommand
	receiver   Receiver
	commandMap map[string]Command
}

type StateCommand struct {
	concreteCommand
	receiver   Receiver
	commandMap map[string]Command
}

type concreteCommand struct {
	commandName        string
	commandDescription string
}

func (c *concreteCommand) GetName() string {
	return c.commandName
}

func (c *concreteCommand) GetDescription() string {
	return fmt.Sprintf("%s -- %s", c.commandName, c.commandDescription)
}

func (mIC *MyIssuesCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return mIC.receiver.myIssues(username)
}

func (cIC *CreateIssueCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return cIC.receiver.createIssue(username, text, cIC.commandMap)
}

func (dC *DoneCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return dC.receiver.done(ctx, username)
}

func (cC *CancelCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return cC.receiver.cancel(username)
}

func (nC *NilCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return nC.receiver.noCommand(username, text, tag, fileID)
}

func (hC *HelpCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return hC.receiver.helpCommand(hC.commandMap)
}

func (sC *StateCommand) execute(ctx context.Context, username, text string, tag string, fileID tgbotapi.FileID) (string, error) {
	return sC.receiver.stateCommand(username, sC.commandMap)
}
