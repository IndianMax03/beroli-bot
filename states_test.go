package main

import (
	"fmt"
	"testing"
)

func generatePositiveCasesGetLocalizedStateDescription() []InputStringWantString {
	cases := []InputStringWantString{
		{
			Input: NIL_STATE,
			Want:  NIL_STATE_RU,
		},
		{
			Input: CREATING_STATE,
			Want:  CREATING_STATE_RU,
		},
		{
			Input: DONE_STATE,
			Want:  DONE_STATE_RU,
		},
		{
			Input: CANCELED_STATE,
			Want:  CANCELED_STATE_RU,
		},
	}

	return cases
}

func generateNegativeCasesGetLocalizedStateDescription() []InputStringWantError {
	cases := []InputStringWantError{
		{
			Input: "",
			Error: ErrStateNotExists,
		},
		{
			Input: RandStringWithSymbolsAndEmojis(10),
			Error: ErrStateNotExists,
		},
	}

	return cases
}

func TestGetLocalizedStateDescriptionPositive(t *testing.T) {
	tests := generatePositiveCasesGetLocalizedStateDescription()
	for _, test := range tests {
		name := fmt.Sprintf("CASE:'%s'->'%s'", test.Input, test.Want)
		t.Run(name, func(t *testing.T) {
			got, err := GetLocalizedStateDescription(test.Input)
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
			_, err := GetLocalizedStateDescription(test.Input)
			if err != test.Error {
				t.Errorf("Got: error is nil, Want error '%v'", test.Error)
			}
		})
	}
}
