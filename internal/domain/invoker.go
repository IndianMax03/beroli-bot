package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrEmptyText      = errors.New("пустой текст")
	ErrUnknownCommand = errors.New("неизвестная команда")
	ErrMultiplyTags   = errors.New("не могу распознать несколько тэгов в одном сообщении")
	ErrUnknownTag     = errors.New("неизвестный тэг")
)

type Invoker struct {
	receiver Receiver
	commands map[string]Command
}

func NewInvoker(r Receiver) *Invoker {
	commandMap := map[string]Command{
		MY_ISSUES_COMMAND: NewMyIssuesCommand(r),
		DONE_COMMAND:      NewDoneCommand(r),
		CANCEL_COMMAND:    NewCancelCommand(r),
		NIL_COMMAND:       NewNilCommand(r),
	}
	commandMap[CREATE_ISSUE_COMMAND] = NewCreateIssueCommand(r, commandMap)
	commandMap[STATE_COMMAND] = NewStateCommand(r, commandMap)
	commandMap[HELP_COMMAND] = NewHelpCommand(r, commandMap)
	return &Invoker{
		receiver: r,
		commands: commandMap,
	}
}

func (i *Invoker) ExecuteCommand(ctx context.Context, username, text string, fileID tgbotapi.FileID) (string, error) {

	cmd, text, tag, err := parseCommand(text, i)
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

	return cmd.Execute(ctx, username, text, tag, fileID)
}

func parseCommand(body string, i *Invoker) (cmd Command, text string, tag string, err error) {
	var b strings.Builder
	for line := range strings.SplitSeq(body, "\n") {
		it := 0
		for word := range strings.SplitSeq(line, " ") {
			if strings.HasPrefix(word, "#") {
				if tag != "" {
					return nil, "", "", ErrMultiplyTags
				}
				trimmedWord := strings.TrimSuffix(word, "\n")
				if _, ok := LocalizedTagsDescriptionMap[trimmedWord]; !ok {
					return nil, "", "", fmt.Errorf("%s: %w", trimmedWord, ErrUnknownTag)
				}
				tag = trimmedWord
			} else if strings.HasPrefix(word, "/") && cmd == nil {
				var ok bool
				cmd, ok = i.commands[word]
				if !ok {
					return nil, "", "", ErrUnknownCommand
				}
			} else {
				if it > 0 {
					b.WriteString(" ")
				}
				if word == "" {
					b.WriteString("\n")
				} else {
					b.WriteString(word)
				}
				it++
			}
		}
	}
	if cmd == nil {
		cmd = i.commands[NIL_COMMAND]
	}
	text = b.String()
	return
}
