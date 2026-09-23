package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type Validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode JSON: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(Validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil

}
