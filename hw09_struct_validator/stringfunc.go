package hw09structvalidator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func IsRussian(s string) bool {
	for _, r := range s {
		// Если символ — буква, но не русская, возвращаем false
		if unicode.IsLetter(r) && !unicode.Is(unicode.Cyrillic, r) {
			return false
		}
	}
	return true
}

func IsChinese(s string) bool {
	for _, r := range s {
		// Если символ — буква, но не китайский иероглиф, возвращаем false
		if unicode.IsLetter(r) && !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}

func IsEng(s string) bool {
	for _, r := range s {
		// Если символ — буква, но не русская, возвращаем false
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			return false
		}
	}
	return true
}

func StrLenIs(s string, l int) bool {
	return len(s) == l
}

func StrLenNotMore(s string, l int) bool {
	return len(s) <= l
}

func StrLenNotLess(s string, l int) bool {
	return len(s) >= l
}

func StrIn(s string, stSet []string) bool {
	for _, st := range stSet {
		if s == st {
			return true
		}
	}
	return false
}

func StrRegexp(s string, reg string) (bool, error) {
	re, err := regexp.Compile(reg)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrRegexTagValidation, err)
	}
	return re.MatchString(s), nil
}

func Email(s string) bool {
	res, _ := StrRegexp(s, "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	return res
}

func PhoneNumber(s string) bool {
	res, _ := StrRegexp(s, "^\\+\\d{11}$")
	return res
}

func StringValidate(s string, valTag string, fName string, valEr *ValidationErrors) {
	rules := strings.Split(valTag, "|")
	for _, rule := range rules {
		switch rule {
		case "russian":
			if !IsRussian(s) {
				valEr.Add(fName,
					fmt.Errorf("(правило: %s), значение=%s: %w", rule, s, ErrStringLangValidation))
			}
		case "chinese":
			if !IsChinese(s) {
				valEr.Add(fName,
					fmt.Errorf("(правило: %s), значение=%s: %w", rule, s, ErrStringLangValidation))
			}
		case "eng":
			if !IsEng(s) {
				valEr.Add(fName,
					fmt.Errorf("(правило: %s), значение=%s: %w", rule, s, ErrStringLangValidation))
			}
		case "e-mail":
			if !Email(s) {
				valEr.Add(fName,
					fmt.Errorf("(правило: %s), значение=%s: %w", rule, s, ErrEmailValidation))
			}
		case "phonenumber":
			if !PhoneNumber(s) {
				valEr.Add(fName,
					fmt.Errorf("(правило: %s), значение=%s: %w", rule, s, ErrPhoneNumberValidation))
			}
		default:
			r := strings.Split(rule, ":")
			if len(r) != 2 {
				valEr.Add(fName, fmt.Errorf("%w: %s, правило: %s", ErrStringTag, valTag, rule))
				return
			}
			switch r[0] {
			case "regexp":
				ok, err := StrRegexp(s, r[1])
				if err != nil {
					valEr.Add(fName, fmt.Errorf("%w: %w", ErrRegexCompile, err))
					return
				}
				if !ok {
					valEr.Add(fName, fmt.Errorf("(правило: %s), значение=%q: %w",
						rule, s, ErrRegexTagValidation))
				}
			case "len":
				l, err := strconv.ParseUint(r[1], 10, 32)
				if err != nil {
					valEr.Add(fName, fmt.Errorf("%w (правило: %s), значение=%s: %w",
						ErrStringTag, rule, r[1], err))
					return
				}
				if !StrLenIs(s, int(l)) {
					valEr.Add(fName,
						fmt.Errorf("%w: длина строки не соответствует заданной (правило: %s), длина строки=%d",
							ErrorStrLen, rule, len(s)))
				}
			case "lennotless":
				l, err := strconv.ParseUint(r[1], 10, 32)
				if err != nil {
					valEr.Add(fName, fmt.Errorf("%w (правило: %s), значение=%s: %w",
						ErrStringTag, rule, r[1], err))
					return
				}
				if !StrLenNotLess(s, int(l)) {
					valEr.Add(fName, fmt.Errorf(""+
						"превышено минимальное ограничение длины строки: (правило: %s), длина строки=%d", rule, len(s)))
				}
			case "lennotmore":
				l, err := strconv.ParseUint(r[1], 10, 32)
				if err != nil {
					valEr.Add(fName, fmt.Errorf("%w (правило: %s), значение=%s: %w",
						ErrStringTag, rule, r[1], err))
					return
				}
				if !StrLenNotMore(s, int(l)) {
					valEr.Add(fName, fmt.Errorf(""+
						"превышено мксиимальное ограничение длины строки: (правило: %s), длина строки=%d",
						rule, len(s)))
				}
			case "in":
				stSet := strings.Split(r[1], ",")
				if len(stSet) < 1 {
					valEr.Add(fName, fmt.Errorf(""+
						"%w: не указан параметр (правило: %s)",
						ErrStringTag, rule))
				}
				if !StrIn(s, stSet) {
					valEr.Add(fName, fmt.Errorf("%w: значение не входит в словарь (правило %s) значение=%q",
						ErrorStringDict, rule, s))
				}
			}
		}
	}
}
