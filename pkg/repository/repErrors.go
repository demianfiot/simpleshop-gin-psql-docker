package repository

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")

	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductNotFound      = errors.New("product not found")

	ErrOrderNotFound = errors.New("order not found")
	ErrAccessDenied  = errors.New("access denied")
)
