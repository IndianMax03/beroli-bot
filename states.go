package main

import "errors"

const (
	NIL_STATE      = "nil"
	CREATING_STATE = "creating"
	DONE_STATE     = "done"
	CANCELED_STATE = "canceled"
)

var ErrStateNotExists = errors.New("неизвестное состояние пользователя")

var localizedStatesDescriptionMap = map[string]string{
	NIL_STATE:      "отсутствует",
	CREATING_STATE: "создание задачи",
	DONE_STATE:     "задача создана",
	CANCELED_STATE: "задача отменена",
}

func GetLocalizedStateDescription(state string) (string, error) {
	if result, ok := localizedStatesDescriptionMap[state]; !ok {
		return "", ErrStateNotExists
	} else {
		return result, nil
	}
}
