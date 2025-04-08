package main

func NewInvoker(r *Receiver) *Invoker {
	commandMap := map[string]Command{
		MY_ISSUES_COMMAND:    NewMyIssuesCommand(r),
		CREATE_ISSUE_COMMAND: NewCreateIssueCommand(r),
		DONE_COMMAND:         NewDoneCommand(r),
		CANCEL_COMMAND:       NewCancelCommand(r),
		NIL_COMMAND:          NewNilCommand(r),
	}
	return &Invoker{
		receiver: r,
		commands: commandMap,
	}
}

type Invoker struct {
	receiver *Receiver
	commands map[string]Command
}

func (i *Invoker) executeCommand(cmd, text string) string {
	command, ok := i.commands[cmd]
	if !ok {
		command = i.commands[NIL_COMMAND]
	}
	return command.execute(text)
}
