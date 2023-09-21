package models

import (
	"errors"
)

var (
	// Error created in case if requested record does not exists
	ErrNoRecord = errors.New("models: no matching record found")
	// Error for when a user tries to login with incorrect email
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Error for when a user tries to signup with an email that is already
	// in a database
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
