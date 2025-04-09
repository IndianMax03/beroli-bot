package main

import "github.com/IndianMax03/yandex-tracker-go-client/v3/model"

const (
	USERNAME_FIELD = "username"
	STATE_FIELD    = "state"
)

type User struct {
	Username string
	State    string
	Issue    *model.IssueCreateRequest
	Issues   []Issue
}

type Issue struct {
	Key  string
	Link string
}
