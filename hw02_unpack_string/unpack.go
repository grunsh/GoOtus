package unpack

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func IsLetter(l rune) bool {
	return unicode.IsLetter(l) ||
		(l >= 'a' && l <= 'z') ||
		l == rune('\n') ||
		l == rune(' ') ||
		l == '_'
}

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	var repeatTimes int
	var tmpString strings.Builder
	prevRune, er := utf8.DecodeRuneInString(s)
	if er == 0 {
		return "", ErrInvalidString
	}
	runeSize := utf8.RuneLen(prevRune)
	if !IsLetter(prevRune) {
		return "", ErrInvalidString
	}
	for _, value := range s[runeSize:] {
		switch {
		case IsLetter(prevRune):
			switch {
			case unicode.IsDigit(value):
				repeatTimes = int(value - '0')
			case IsLetter(value):
				repeatTimes = 1
			default:
				return "", ErrInvalidString
			}
			tmpString.WriteString(strings.Repeat(string(prevRune), repeatTimes))
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
