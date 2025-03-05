package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func CheckAllErrors(requestError []error, valErrors ValidationErrors) (allFounded bool) {
	for _, err := range requestError {
		allFounded = false
		for _, e := range valErrors {
			if errors.Is(e.Err, err) {
				allFounded = true
				break
			}
		}
	}
	return
}

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name         string
		in           interface{}
		expectedErrs []error
	}{
		{
			name: "Корректный номер телефона",
			in: struct {
				phonenumber string `validate:"phonenumber"`
			}{
				phonenumber: "+79991234567",
			},
			expectedErrs: nil,
		},
		{
			name: "Не корректный номер телефона",
			in: struct {
				phonenumber string `validate:"phonenumber"`
			}{
				phonenumber: "-79991234567",
			},
			expectedErrs: []error{ErrPhoneNumberValidation},
		},
		{
			name: "Корректный юзер",
			in: User{
				ID:     "012345678901234567890123456789123456",
				Name:   "вася",
				Age:    19,
				Email:  "non@mail.com",
				Role:   "admin",
				Phones: []string{"79067240163", "79991234567"},
				meta:   []byte{},
			},
			expectedErrs: nil,
		},
		{
			name: "Сломанный полностью юзер",
			in: User{
				ID:     "01234567890123456789012345678912345",
				Name:   "вася",
				Age:    10,
				Email:  "()non@mail.com",
				Role:   "admin_",
				Phones: []string{"9067240163", "9991234567"},
				meta:   []byte{},
			},
			expectedErrs: []error{
				ErrorStrLen,
				ErrorIntMin,
				ErrRegexTagValidation,
				ErrorStrLen,
			},
		},
		{
			name: "Китайский правильный",
			in: struct {
				chinaWord string `validate:"chinese"`
			}{
				chinaWord: "嘿混蛋",
			},
			expectedErrs: nil,
		},
		{
			name: "Китайский сломанный",
			in: struct {
				chinaWord string `validate:"chinese"`
			}{
				chinaWord: "嘿混蛋апчхи",
			},
			expectedErrs: []error{
				ErrStringLangValidation,
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			//			tt := tt
			t.Parallel()

			ve := Validate(tt.in)
			// fmt.Println("================ "+tt.name+"   ", ve)
			// fmt.Println("================ "+tt.name+"   ", tt.expectedErrs)
			if tt.expectedErrs == nil {
				require.Zero(t, len(ve))
			} else {
				require.True(t, CheckAllErrors(tt.expectedErrs, ve))
			}
		})
	}
}
