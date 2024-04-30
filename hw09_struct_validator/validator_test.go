package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "12345678901234567890123456789012345",
				Name:   "John",
				Age:    15,
				Email:  "!123@test.com",
				Role:   "admin2",
				Phones: []string{"123456789012"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrStringNotLen},
				{Field: "Age", Err: ErrIntUnderMin},
				{Field: "Email", Err: ErrStringNotRegexp},
				{Field: "Role", Err: ErrStringNotInSet},
				{Field: "Phones", Err: ErrStringNotLen},
			},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    90,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: ErrIntOverMax},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrStringNotLen},
			},
		},
		{
			in: Token{
				Header:    []byte("1234567890"),
				Payload:   []byte("12345678901234567890"),
				Signature: []byte("12345678901234567890123456789012"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
				Body: "Created",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrIntNotInSet},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			errs := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, errs)
			} else {
				var valErrs ValidationErrors
				if errors.As(errs, &valErrs) {
					var expectedErrs ValidationErrors
					require.ErrorAs(t, tt.expectedErr, &expectedErrs)
					require.Equal(t, len(expectedErrs), len(valErrs),
						fmt.Sprintf(
							"expected err: %s\ngot err: %s\n",
							expectedErrs, valErrs,
						),
					)
					for i, err := range valErrs {
						require.ErrorIs(t, err, expectedErrs[i])
					}
				} else {
					require.ErrorIs(t, errs, tt.expectedErr)
				}
			}
			_ = tt
		})
	}
}
