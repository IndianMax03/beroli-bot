package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyText      = errors.New("can't parse command because text is empty")
	ErrUnknownCommand = errors.New("unknown command")
)

func NewInvoker(r Receiver) *Invoker {
	commandMap := map[string]Command{
		MY_ISSUES_COMMAND:    NewMyIssuesCommand(r),
		CREATE_ISSUE_COMMAND: NewCreateIssueCommand(r),
		DONE_COMMAND:         NewDoneCommand(r),
		CANCEL_COMMAND:       NewCancelCommand(r),
		NIL_COMMAND:          NewNilCommand(r),
	}
	commandMap[STATE_COMMAND] = NewStateCommand(r, commandMap)
	commandMap[HELP_COMMAND] = NewHelpCommand(r, commandMap)
	return &Invoker{
		receiver: r,
		commands: commandMap,
	}
}

type Invoker struct {
	receiver Receiver
	commands map[string]Command
}

func (i *Invoker) executeCommand(username, text string) (string, error) {

	cmd, text, tags, err := parseCommand(text, i)
	if err != nil {
		return "", err
	}

	err = i.receiver.ValidateAndInitUser(username)
	if err != nil {
		return "", err
	}

	err = i.receiver.ValidateState(username, cmd.GetName())
	if err != nil {
		return "", err
	}

	return cmd.execute(username, text, tags)
}

func parseCommand(body string, i *Invoker) (cmd Command, text string, tags []string, err error) {
	if body == "" {
		return nil, "", nil, ErrEmptyText
	}
	var b strings.Builder
	for _, word := range strings.Split(body, " ") {
		if strings.HasPrefix(word, "#") {
			tags = append(tags, word)
		} else if strings.HasPrefix(word, "/") && cmd == nil {
			var ok bool
			cmd, ok = i.commands[word]
			if !ok {
				return nil, "", nil, ErrUnknownCommand
			}
		} else {
			b.WriteString(fmt.Sprintf("%s ", word))
		}
	}
	if cmd == nil {
		cmd = i.commands[NIL_COMMAND]
	}
	text = b.String()
	return
}
