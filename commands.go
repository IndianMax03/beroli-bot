package main

import "fmt"

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

func NewCreateIssueCommand(r Receiver) *CreateIssueCommand {
	return &CreateIssueCommand{
		concreteCommand: concreteCommand{
			commandName:        CREATE_ISSUE_COMMAND,
			commandDescription: CREATE_ISSUE_COMMAND_DESCRIPTION,
		},
		receiver: r,
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

type Command interface {
	execute(string, string, []string) (string, error)
	GetName() string
	GetDescription() string
}

type MyIssuesCommand struct {
	concreteCommand
	receiver Receiver
}

type CreateIssueCommand struct {
	concreteCommand
	receiver Receiver
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

func (mIC *MyIssuesCommand) execute(username, text string, tags []string) (string, error) {
	return mIC.receiver.myIssues()
}

func (cIC *CreateIssueCommand) execute(username, text string, tags []string) (string, error) {
	return cIC.receiver.createIssue(username, text)
}

func (dC *DoneCommand) execute(username, text string, tags []string) (string, error) {
	return dC.receiver.done()
}

func (cC *CancelCommand) execute(username, text string, tags []string) (string, error) {
	return cC.receiver.cancel()
}

func (nC *NilCommand) execute(username, text string, tags []string) (string, error) {
	return nC.receiver.noCommand(text)
}
