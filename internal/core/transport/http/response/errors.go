package core_http_response

import (
	"errors"
	"net/http"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

type ErrorResponse struct {
	TextCode string `json:"text_code" example:"invalid_argument"`
	Message  string `json:"message" example:"the provided argument is invalid"`
}

func StatusFromError(err error) int {
	switch {
	case errors.Is(err, core_errors.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, core_errors.ErrInvalidArgument):
		return http.StatusBadRequest
	case errors.Is(err, core_errors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, core_errors.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, core_errors.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func CodeFromError(err error) string {
	switch {
	case errors.Is(err, core_errors.ErrUsernameTaken):
		return "username_taken"
	case errors.Is(err, core_errors.ErrEmailTaken):
		return "email_taken"
	case errors.Is(err, core_errors.ErrConflict):
		return "conflict"
	case errors.Is(err, core_errors.ErrInvalidArgument):
		return "invalid_argument"
	case errors.Is(err, core_errors.ErrNotFound):
		return "not_found"
	case errors.Is(err, core_errors.ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, core_errors.ErrForbidden):
		return "forbidden"
	default:
		return "internal"
	}
}
