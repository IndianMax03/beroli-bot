package main

import "strings"

const MAX_STRING_LENGHT = 30
const MAX_ARRAY_LENGHT = 2

func CutString(s string) string {
	runes := []rune(s)
	var b strings.Builder
	for i, r := range runes {
		if i == MAX_STRING_LENGHT {
			b.WriteString("...")
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}

func CutArrayOfString(ss []string) string {
	result := []string{}
	for i, s := range ss {
		if i == MAX_ARRAY_LENGHT {
			result = append(result, "...")
			break
		}
		result = append(result, CutString(s))
	}
	return strings.Join(result[:], ", ")
}
