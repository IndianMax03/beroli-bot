package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IndianMax03/beroli-bot/internal/util"
)

const (
	CUT_STRING_CASES_COUNT = 48
	CUT_ARRAY_CASES_COUNT  = 48
)

var (
	cutStringLengthBoundaryValues = []int{
		0,
		1,
		util.MAX_STRING_LENGHT - 1,
		util.MAX_STRING_LENGHT,
		util.MAX_STRING_LENGHT + 1,
		util.MAX_STRING_LENGHT + util.MAX_STRING_LENGHT,
	}
	cutArrayLengthBoundaryValues = []int{
		util.MAX_ARRAY_LENGHT + util.MAX_ARRAY_LENGHT,
		util.MAX_ARRAY_LENGHT + 1,
		util.MAX_ARRAY_LENGHT,
		util.MAX_ARRAY_LENGHT - 1,
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
		if len(inputRunes) > util.MAX_STRING_LENGHT {
			inputRunes = inputRunes[:util.MAX_STRING_LENGHT]
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
			if i == util.MAX_ARRAY_LENGHT {
				cutInput = append(cutInput, "...")
				break
			}
			cutInput = append(cutInput, util.CutString(s))
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
			got := util.CutString(test.Input)
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
			got := util.CutArrayOfString(test.Input)
			if got != test.Want {
				t.Errorf("Got: '%s', Want '%s'", got, test.Want)
			}
		})
	}
}
