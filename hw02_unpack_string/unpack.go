package unpack

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func IsLetter(l rune) bool {
	return (l >= 'A' && l <= 'Z') || (l >= 'a' && l <= 'z') || l == rune('\n')
}

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	var RepeatTimes int
	var tmpString strings.Builder
	prevRune, _ := utf8.DecodeRuneInString(s)
	runeSize := utf8.RuneLen(prevRune)
	if !IsLetter(prevRune) {
		return "", ErrInvalidString
	}
	for _, value := range s[runeSize:] {
		switch {
		case IsLetter(prevRune):
			switch {
			case unicode.IsDigit(value):
				RepeatTimes = int(value - '0')
			case IsLetter(value):
				RepeatTimes = 1
			}
			tmpString.WriteString(strings.Repeat(string(prevRune), RepeatTimes))
		case !unicode.IsDigit(prevRune):
			return "", ErrInvalidString
		case unicode.IsDigit(value):
			return "", ErrInvalidString
		}
		prevRune = value
	}
	lastRune, _ := utf8.DecodeLastRuneInString(s)
	switch {
	case IsLetter(lastRune):
		tmpString.WriteString(string(lastRune))
	case unicode.IsDigit(lastRune):
	default:
		return "", ErrInvalidString
	}
	return tmpString.String(), nil
}
