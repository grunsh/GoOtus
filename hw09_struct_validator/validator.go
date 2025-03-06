package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/constraints"
)

var (
	ErrorNotStruct           = errors.New("подсуноли что-то, но не структуру")
	ErrorValidationTag       = errors.New("не корретный тэг валиадации")
	ErrorValidationRange     = errors.New("диапазон валиадции не соответствует типу")
	ErrStringTag             = errors.New("ошибка тега валиадации строки")
	ErrStringLangValidation  = errors.New("строка содержит символы недопустимого языка")
	ErrEmailValidation       = errors.New("ошибка в формате адреса электронной почты")
	ErrRegexCompile          = errors.New("ошибка компиляции регулярного выражения тега валидации")
	ErrRegexTagValidation    = errors.New("ошибка проверки строки с помощью шаблона")
	ErrPhoneNumberValidation = errors.New("ошибка в номере телефона")
	ErrorStrLen              = errors.New("ошибка длины строки")
	ErrorStringDict          = errors.New("нет в словаре")
	ErrorIntMin              = errors.New("ниже минимального")
	ErrorIntMax              = errors.New("больше  максимального")
	ErrorIntNotInSet         = errors.New("целое не входит в множество")
)

type ValidationError struct {
	Field string
	Err   error
}

type AllInt constraints.Integer

func IsStruct(v interface{}) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr { // Если v является указателем, получаем значение, на которое он указывает
		val = val.Elem()
	}
	return val.Kind() == reflect.Struct
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errorMessages := make([]string, 0)
	for _, err := range v {
		errorMessages = append(errorMessages, fmt.Sprintf("Field: %s, Error: %s", err.Field, err.Err.Error()))
	}
	return strings.Join(errorMessages, "\n")
}

// Метод для добавления новой ошибки в слайс.
func (v *ValidationErrors) Add(field string, err error) {
	*v = append(*v, ValidationError{Field: field, Err: err})
}

func Validate(v interface{}) ValidationErrors {
	if !IsStruct(v) {
		return []ValidationError{{"No fields", ErrorNotStruct}}
	}

	// Если вдруг указатель, то возьмём значение по указателю
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	vErr := ValidationErrors{}
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		validationTag := field.Tag.Get("validate")
		if validationTag == "" { // Пропускаем поля, для которых тэг валидации не указан.
			continue
		}
		fieldName := field.Name
		fieldType := field.Type.Kind().String()
		fieldValue := val.Field(i)
		switch fieldValue.Kind() {
		case reflect.String:
			StringValidate(fieldValue.String(), validationTag, fieldName, &vErr)
		case reflect.Slice:
			elemType := fieldValue.Type().Elem()
			elemStrType := elemType.Kind().String()
			switch elemType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				for j := 0; j < fieldValue.Len(); j++ {
					IntValidate(fieldValue.Index(j).Int(), validationTag, elemStrType,
						fmt.Sprintf("%s[%d]", fieldName, j), &vErr)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				for j := 0; j < fieldValue.Len(); j++ {
					IntValidate(fieldValue.Index(j).Uint(), validationTag, elemStrType,
						fmt.Sprintf("%s[%d]", fieldName, j), &vErr)
				}
			case reflect.String:
				for j := 0; j < fieldValue.Len(); j++ {
					StringValidate(fieldValue.Index(j).String(), validationTag,
						fmt.Sprintf("%s[%d]", fieldName, j), &vErr)
				}
			default:
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			IntValidate(fieldValue.Int(), validationTag, fieldType, fieldName, &vErr)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			IntValidate(fieldValue.Uint(), validationTag, fieldType, fieldName, &vErr)
		default:
		}
	}

	return vErr
}
