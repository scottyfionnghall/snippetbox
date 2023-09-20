package validator

import (
	"strings"
	"unicode/utf8"
)

// Define a new Validator type wich contains a mp of validation errors for our
// form fields.

type Validator struct {
	FieldErrors map[string]string
}

// Valid() return true if the FieldErrors map doesn't contain
// any entries
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map
// (so long as no entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// We need to initilize the map first, if it isn't already
	// initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exsits := v.FieldErrors[key]; !exsits {
		v.FieldErrors[key] = message
	}
}

// CheckField() adds an error mesage to the FieldErrors map
// only if a validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank() returns true if value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermitedInt() returns true if a value is in a list of permitted integers.
func PermitedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
