package test

import (
	"fmt"
	"testing"

	"github.com/IndianMax03/beroli-bot/internal/domain"
)

func generatePositiveCasesGetLocalizedStateDescription() []InputStringWantString {
	cases := []InputStringWantString{
		{
			Input: domain.NIL_STATE,
			Want:  domain.NIL_STATE_RU,
		},
		{
			Input: domain.CREATING_STATE,
			Want:  domain.CREATING_STATE_RU,
		},
		{
			Input: domain.DONE_STATE,
			Want:  domain.DONE_STATE_RU,
		},
		{
			Input: domain.CANCELED_STATE,
			Want:  domain.CANCELED_STATE_RU,
		},
	}

	return cases
}

func generateNegativeCasesGetLocalizedStateDescription() []InputStringWantError {
	cases := []InputStringWantError{
		{
			Input: "",
			Error: domain.ErrStateNotExists,
		},
		{
			Input: RandStringWithSymbolsAndEmojis(10),
			Error: domain.ErrStateNotExists,
		},
	}

	return cases
}

func TestGetLocalizedStateDescriptionPositive(t *testing.T) {
	tests := generatePositiveCasesGetLocalizedStateDescription()
	for _, test := range tests {
		name := fmt.Sprintf("CASE:'%s'->'%s'", test.Input, test.Want)
		t.Run(name, func(t *testing.T) {
			got, err := domain.GetLocalizedStateDescription(test.Input)
			if err != nil {
				t.Errorf("Got error: '%v', Want '%s'", err, test.Want)
			} else if got != test.Want {
				t.Errorf("Got: '%s', Want '%s'", got, test.Want)
			}
		})
	}
}

func TestGetLocalizedStateDescriptionNegative(t *testing.T) {
	tests := generateNegativeCasesGetLocalizedStateDescription()
	for _, test := range tests {
		name := fmt.Sprintf("CASE:'%s'->'%v'", test.Input, test.Error)
		t.Run(name, func(t *testing.T) {
			_, err := domain.GetLocalizedStateDescription(test.Input)
			if err != test.Error {
				t.Errorf("Got: error is nil, Want error '%v'", test.Error)
			}
		})
	}
}
