package core_errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrConflict        = errors.New("conflict")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")

	// Частные случаи конфликта: оборачивают ErrConflict, поэтому статус остаётся 409,
	// а клиент по коду отличает «имя занято» от «email занят».
	ErrUsernameTaken = fmt.Errorf("username taken: %w", ErrConflict)
	ErrEmailTaken    = fmt.Errorf("email taken: %w", ErrConflict)
)
