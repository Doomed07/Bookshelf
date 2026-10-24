package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	"go.uber.org/zap"
)

// zap.NewNop() — логгер, который ничего не пишет: в тесте нам не нужны
// ни файлы логов, ни вывод в консоль
func newTestLogger() *core_logger.Logger {
	return &core_logger.Logger{Logger: zap.NewNop()}
}

func TestHTTPResponseHandler_ErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "invalid argument", err: core_errors.ErrInvalidArgument, wantStatus: http.StatusBadRequest, wantCode: "invalid_argument"},
		{name: "not found", err: core_errors.ErrNotFound, wantStatus: http.StatusNotFound, wantCode: "not_found"},
		{name: "conflict", err: core_errors.ErrConflict, wantStatus: http.StatusConflict, wantCode: "conflict"},
		{name: "username taken", err: core_errors.ErrUsernameTaken, wantStatus: http.StatusConflict, wantCode: "username_taken"},
		{name: "email taken", err: core_errors.ErrEmailTaken, wantStatus: http.StatusConflict, wantCode: "email_taken"},
		{name: "unauthorized", err: core_errors.ErrUnauthorized, wantStatus: http.StatusUnauthorized, wantCode: "unauthorized"},
		{name: "forbidden", err: core_errors.ErrForbidden, wantStatus: http.StatusForbidden, wantCode: "forbidden"},
		{
			name:       "wrapped not found",
			err:        fmt.Errorf("get user: %w", fmt.Errorf("repo: %w", core_errors.ErrNotFound)),
			wantStatus: http.StatusNotFound,
			wantCode:   "not_found",
		},
		{name: "unknown error", err: errors.New("db is down"), wantStatus: http.StatusInternalServerError, wantCode: "internal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// httptest.ResponseRecorder запоминает статус и тело ответа
			rec := httptest.NewRecorder()
			h := NewHTTPResponseHandler(newTestLogger(), rec)

			h.ErrorResponse(tt.err, "failed to do something")

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var body ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode response body: %v", err)
			}
			if body.TextCode != tt.wantCode || body.Message != "failed to do something" {
				t.Errorf("body = %+v, want text_code=%q message=%q", body, tt.wantCode, "failed to do something")
			}
		})
	}
}

func TestHTTPResponseHandler_PanicResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	h := NewHTTPResponseHandler(newTestLogger(), rec)

	h.PanicResponse("boom", "panic in handler")

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestHTTPResponseHandler_ContentType(t *testing.T) {
	tests := []struct {
		name string
		call func(h *HTTPResponseHandler)
	}{
		{"JSONResponse", func(h *HTTPResponseHandler) { h.JSONResponse(http.StatusOK, map[string]string{"a": "b"}) }},
		{"ErrorResponse", func(h *HTTPResponseHandler) { h.ErrorResponse(core_errors.ErrNotFound, "msg") }},
		{"PanicResponse", func(h *HTTPResponseHandler) { h.PanicResponse("boom", "msg") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.call(NewHTTPResponseHandler(newTestLogger(), rec))

			// Result() — это то, что получил бы клиент: заголовки зафиксированы в момент WriteHeader
			if got := rec.Result().Header.Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}

			if got := rec.Result().Header.Get("X-Content-Type-Options"); got != "nosniff" {
				t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
			}
		})
	}
}

// Текст исходной ошибки (имена таблиц, детали БД, текст паники) не должен попадать клиенту: он остаётся только в логах.
func TestHTTPResponseHandler_NoErrorTextLeak(t *testing.T) {
	const secret = "secret db detail"

	tests := []struct {
		name string
		call func(h *HTTPResponseHandler)
	}{
		{"known error", func(h *HTTPResponseHandler) {
			h.ErrorResponse(fmt.Errorf("%s: %w", secret, core_errors.ErrNotFound), "msg")
		}},
		{"unknown error", func(h *HTTPResponseHandler) { h.ErrorResponse(errors.New(secret), "msg") }},
		{"panic", func(h *HTTPResponseHandler) { h.PanicResponse(secret, "msg") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tt.call(NewHTTPResponseHandler(newTestLogger(), rec))

			if strings.Contains(rec.Body.String(), secret) {
				t.Errorf("response body leaks error text: %s", rec.Body.String())
			}
		})
	}
}
