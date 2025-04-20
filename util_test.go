package main

import (
	"fmt"
	"strings"
	"testing"
)

const (
	CUT_STRING_CASES_COUNT = 48
	CUT_ARRAY_CASES_COUNT  = 48
)

var (
	cutStringLengthBoundaryValues = []int{
		0,
		1,
		MAX_STRING_LENGHT - 1,
		MAX_STRING_LENGHT,
		MAX_STRING_LENGHT + 1,
		MAX_STRING_LENGHT + MAX_STRING_LENGHT,
	}
	cutArrayLengthBoundaryValues = []int{
		MAX_ARRAY_LENGHT + MAX_ARRAY_LENGHT,
		MAX_ARRAY_LENGHT + 1,
		MAX_ARRAY_LENGHT,
		MAX_ARRAY_LENGHT - 1,
		1,
		0,
	}
)

func genrateCasesCutString() []InputStringWantString {
	var cases []InputStringWantString
	caseNum := 0
	for range CUT_STRING_CASES_COUNT {
		input := RandStringWithSymbolsAndEmojis(cutStringLengthBoundaryValues[caseNum])
		inputRunes := []rune(input)
		want := input
		if len(inputRunes) > MAX_STRING_LENGHT {
			inputRunes = inputRunes[:MAX_STRING_LENGHT]
			want = string(inputRunes) + "..."
		}
		cases = append(cases, InputStringWantString{
			Input: input,
			Want:  want,
		})
		if caseNum == len(cutStringLengthBoundaryValues)-1 {
			caseNum = 0
		} else {
			caseNum++
		}
	}
	return cases
}

func genrateCasesCutArrayOfString() []InputArrayOfStringWantString {
	var cases []InputArrayOfStringWantString
	caseNum := 0
	for range CUT_ARRAY_CASES_COUNT {
		arrayLength := cutArrayLengthBoundaryValues[caseNum]
		input := []string{}
		strCaseNum := 0
		for range arrayLength {
			inputStr := RandStringWithSymbolsAndEmojis(cutStringLengthBoundaryValues[strCaseNum])

			input = append(input, inputStr)

			if strCaseNum == len(cutStringLengthBoundaryValues)-1 {
				strCaseNum = 0
			} else {
				strCaseNum++
			}
		}
		cutInput := []string{}
		for i, s := range input {
			if i == MAX_ARRAY_LENGHT {
				cutInput = append(cutInput, "...")
				break
			}
			cutInput = append(cutInput, CutString(s))
		}
		want := strings.Join(cutInput[:], ", ")
		cases = append(cases, InputArrayOfStringWantString{
			Input: input,
			Want:  want,
		})
		if caseNum == len(cutArrayLengthBoundaryValues)-1 {
			caseNum = 0
		} else {
			caseNum++
		}
	}
	return cases
}

func TestCutString(t *testing.T) {
	tests := genrateCasesCutString()
	for _, test := range tests {
		name := fmt.Sprintf("CASE:'%s'->'%s'", test.Input, test.Want)
		t.Run(name, func(t *testing.T) {
			got := CutString(test.Input)
			if got != test.Want {
				t.Errorf("Got: '%s', Want '%s'", got, test.Want)
			}
		})
	}
}

func TestCutArrayOfString(t *testing.T) {
	tests := genrateCasesCutArrayOfString()
	for _, test := range tests {
		name := fmt.Sprintf("CASE:len=%v->'%s'", len(test.Input), test.Want)
		t.Run(name, func(t *testing.T) {
			got := CutArrayOfString(test.Input)
			if got != test.Want {
				t.Errorf("Got: '%s', Want '%s'", got, test.Want)
			}
		})
	}
}
