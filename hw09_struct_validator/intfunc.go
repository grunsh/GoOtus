package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func MinTag[T AllInt](i T, param any, signed bool) (bool, error) {
	if signed {
		if int64(i) < param.(int64) {
			return false, fmt.Errorf("%w: значение поля %d меньше минимального %d",
				ErrorIntMin, i, param.(int64))
		}
	} else {
		if uint64(i) < param.(uint64) {
			return false, fmt.Errorf("%w: значение поля %d меньше минимального %d",
				ErrorIntMin, i, param.(int64))
		}
	}
	return true, nil
}

func MaxTag[T AllInt](i T, param any, signed bool) (bool, error) {
	if signed {
		if int64(i) > param.(int64) {
			return false, fmt.Errorf("%w: значение поля %d больше максимального %d",
				ErrorIntMax, i, param.(int64))
		}
	} else {
		if uint64(i) > param.(uint64) {
			return false, fmt.Errorf("%w: значение поля %d больше максимального %d",
				ErrorIntMax, i, param.(int64))
		}
	}
	return true, nil
}

func InSet[T AllInt](i T, param []any, signed bool) bool {
	if signed {
		for _, v := range param {
			if int64(i) == v.(int64) {
				return true
			}
		}
	} else {
		for _, v := range param {
			if uint64(i) == v.(uint64) {
				return true
			}
		}
	}
	return false
}

func IntValidation[T AllInt](i T, validationTag string, strType string) (b bool, e error) {
	var signed bool

	var Params []interface{}

	b = true
	e = nil
	putErr := func(e error, put error) error {
		if e != nil {
			return fmt.Errorf("%w: %w", put, e)
		}
		return put
	}

	rules := strings.Split(validationTag, "|")
	for _, rule := range rules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			return false, fmt.Errorf("%w: %s, правило: %s", ErrorValidationTag, validationTag, rule)
		}
		v := strings.Split(r[1], ",")
		if len(v) == 0 {
			return false, fmt.Errorf("не указан параметр валидации для %s в теге %s", r[0], validationTag)
		}
		paramsCount := len(v)
		if strType[:1] == "u" {
			signed = false
		} else {
			signed = true
		}
		Params = make([]interface{}, paramsCount)
		for i := 0; i < paramsCount; i++ {
			val, intErr := ParseAnyInt(v[i], strType)
			if intErr != nil {
				return false, intErr
			}
			Params[i] = val
		}

		switch r[0] {
		case "min":
			ok, er := MinTag(i, Params[0], signed)
			if !ok {
				e = putErr(e, fmt.Errorf("(правило: %s) %w", rule, er))
				// e = fmt.Errorf("%w: (правило: %s)%w:", e, rule, er)
				b = false
			}
		case "max":
			ok, er := MaxTag(i, Params[0], signed)
			if !ok {
				e = putErr(e, fmt.Errorf("(правило: %s)%w", rule, er))
				b = false
			}
		case "in":
			if !InSet(i, Params, signed) {
				e = putErr(e, fmt.Errorf("%w: (правило: %s) значние поля %d", ErrorIntNotInSet, rule, i))
				b = false
			}
		case "out":
			if InSet(i, Params, signed) {
				e = putErr(e, fmt.Errorf("(правило: %s) Значние поля %d входит в множество исключений",
					rule, i))
				b = false
			}
		default:
			e = putErr(e, fmt.Errorf("%w: %s", ErrorValidationTag, rule))
		}
	}
	return
}

func ParseAnyInt(value string, t string) (any, error) {
	var (
		Pars          func(val string, bitSize int) (any, error)
		BitSizeString string
	)
	switch t[:1] {
	case "u":
		if t == "uint" { // Принимаем, что на 32 разрядных системах уже никто не работает
			BitSizeString = "64"
		} else {
			BitSizeString = t[4:] // 8, 16 или 32
		}
		Pars = func(val string, bitSize int) (any, error) {
			return strconv.ParseUint(val, 10, bitSize)
		}
	case "i":
		if t == "int" { // Принимаем, что на 32 разрядных системах уже никто не работает
			BitSizeString = "64"
		} else {
			BitSizeString = t[3:] // 8, 16 или 32
		}
		Pars = func(val string, bitSize int) (any, error) {
			return strconv.ParseInt(val, 10, bitSize)
		}
	default:
		return 0, fmt.Errorf("не распознан тип: %s", t)
	}
	bitSize, bitErr := strconv.Atoi(BitSizeString)
	if bitErr != nil {
		if errors.Is(bitErr, strconv.ErrRange) {
			return 0, fmt.Errorf("%w: %w", ErrorValidationRange, bitErr)
		}
		return 0, bitErr
	}
	retVal, er := Pars(value, bitSize)
	if errors.Is(er, strconv.ErrRange) {
		return 0, fmt.Errorf("%w: %w", ErrorValidationRange, er)
	}
	return retVal, er
}

func IntValidate[T AllInt](i T, validationTag string, strType string, fieldName string, valErr *ValidationErrors) {
	ok, err := IntValidation(i, validationTag, strType)
	if !ok {
		valErr.Add(fieldName, err)
	}
}
