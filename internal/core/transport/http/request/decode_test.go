package core_http_request

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
)

type tagValidatedRequest struct {
	BookID int `json:"book_id" validate:"required,min=1"`
}

type selfValidatedRequest struct {
	Name string `json:"name"`
}

func (r selfValidatedRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is empty")
	}
	return nil
}

func TestDecodeAndValidateRequest(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		dest    any
		wantErr bool
	}{
		{name: "valid body, validator tags", body: `{"book_id": 1}`, dest: &tagValidatedRequest{}},
		{name: "invalid JSON", body: `{"book_id": `, dest: &tagValidatedRequest{}, wantErr: true},
		{name: "wrong type", body: `{"book_id": "one"}`, dest: &tagValidatedRequest{}, wantErr: true},
		{name: "validator tag fails", body: `{"book_id": 0}`, dest: &tagValidatedRequest{}, wantErr: true},
		{name: "missing required field", body: `{}`, dest: &tagValidatedRequest{}, wantErr: true},
		// если тип реализует Validatable, вызывается его Validate(), а не теги
		{name: "valid body, own Validate()", body: `{"name": "kek"}`, dest: &selfValidatedRequest{}},
		{name: "own Validate() fails", body: `{"name": ""}`, dest: &selfValidatedRequest{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			gotErr := DecodeAndValidateRequest(r, tt.dest)

			if tt.wantErr {
				if !errors.Is(gotErr, core_errors.ErrInvalidArgument) {
					t.Errorf("DecodeAndValidateRequest() error = %v, want wrapped core_errors.ErrInvalidArgument", gotErr)
				}
				return
			}

			if gotErr != nil {
				t.Errorf("DecodeAndValidateRequest() failed: %v", gotErr)
			}
		})
	}
}
