package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		in          interface{}
		expectedErr error
	}{
		{User{
			ID:     "",
			Name:   "",
			Age:    70,
			Email:  "",
			Role:   "",
			Phones: []string{""},
			meta:   json.RawMessage(""),
		}, ValidationErrors([]ValidationError{
			{Field: "ID", Err: ErrLength},
			{Field: "Age", Err: ErrMax},
			{Field: "Email", Err: ErrRegexp},
			{Field: "Role", Err: ErrNotIn},
			{Field: "Phones", Err: ErrLength},
		})},
		{User{
			ID:     "11111",
			Name:   "Sergio",
			Age:    12,
			Email:  "test@example.com",
			Role:   "admin",
			Phones: []string{"1239997788"},
		}, ValidationErrors([]ValidationError{
			{Field: "ID", Err: ErrLength},
			{Field: "Age", Err: ErrMin},
			{Field: "Phones", Err: ErrLength},
		})},
		{App{
			Version: "1234",
		}, ValidationErrors{ValidationError{
			"Version", ErrLength,
		}}},
		{App{Version: "11111"}, nil},
		{Response{Code: 200, Body: ""}, nil},
		{Response{
			Code: 310,
			Body: "",
		}, ValidationErrors{ValidationError{
			Field: "Code",
			Err:   ErrNotIn,
		}}},
		{Token{}, nil},
		{nil, ValidationErrors{ValidationError{
			"", ErrInvalidStruct,
		}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, &ValidationErrors{})
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}
