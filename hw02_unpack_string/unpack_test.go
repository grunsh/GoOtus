package unpack

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "a1a1a1b1", expected: "aaab"},
		{input: "abcdcD2S2", expected: "abcdcDDSS"},
		{input: "a3a3", expected: "aaaaaa"},
		{input: "屁_股3", expected: "屁_股股股"},
		{input: "вас2ал", expected: "вассал"},
		{input: "Р1о1с2и1я1", expected: "Россия"},
		{input: "в Отусе пишут н0е0 0чёткие ТЗ на ДЗ", expected: "в Отусе пишут чёткие ТЗ на ДЗ"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "\n\n12", "2", "12", "0_", "0sd3_"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
