package main

import "errors"

const (
	NIL_STATE         = "nil"
	CREATING_STATE    = "creating"
	DONE_STATE        = "done"
	CANCELED_STATE    = "canceled"
	NIL_STATE_RU      = "отсутствует"
	CREATING_STATE_RU = "создание задачи"
	DONE_STATE_RU     = "задача создана"
	CANCELED_STATE_RU = "задача отменена"
)

var ErrStateNotExists = errors.New("неизвестное состояние пользователя")

var localizedStatesDescriptionMap = map[string]string{
	NIL_STATE:      NIL_STATE_RU,
	CREATING_STATE: CREATING_STATE_RU,
	DONE_STATE:     DONE_STATE_RU,
	CANCELED_STATE: CANCELED_STATE_RU,
}

func GetLocalizedStateDescription(state string) (string, error) {
	if result, ok := localizedStatesDescriptionMap[state]; !ok {
		return "", ErrStateNotExists
	} else {
		return result, nil
	}
}
