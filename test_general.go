package main

import (
	mrand "math/rand"
	"strings"
)

const (
	SEED     = 31415926
	ALPHABET = "фбвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SYMBOLS  = " +-='\"!@#$%^&*()[]{}_~`"
	EMOJIS   = "😊😇😂😡😤😭😀"
)

var (
	Rand                  *mrand.Rand
	Alphabet              = []rune(ALPHABET)
	AlphabetSymbols       = []rune(ALPHABET + SYMBOLS)
	AlphabetSymbolsEmojis = []rune(ALPHABET + SYMBOLS + EMOJIS)
)

type InputStringWantString struct {
	Input string
	Want  string
}

type InputStringWantError struct {
	Input string
	Error error
}

type InputArrayOfStringWantString struct {
	Input []string
	Want  string
}

func init() {
	Rand = mrand.New(mrand.NewSource(SEED))
}

func RandStringWithSymbolsAndEmojis(length int) string {
	var b strings.Builder
	b.Grow(length)
	for range length {
		b.WriteRune(AlphabetSymbolsEmojis[Rand.Intn(len(AlphabetSymbolsEmojis))])
	}
	return b.String()
}
