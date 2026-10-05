package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_middleware "github.com/Doomed07/Bookshelf/internal/core/transport/http/middleware"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
	"go.uber.org/zap"
)

type fakeBooks struct{ err error }

func (f fakeBooks) GetBook(ctx context.Context, id int) (core_domain.Book, error) {
	return core_domain.Book{ID: id}, f.err
}

type fakeUsers struct{ err error }

func (f fakeUsers) GetUser(ctx context.Context, id int) (core_domain.User, error) {
	return core_domain.User{ID: id}, f.err
}

func TestRoutes(t *testing.T) {

	tests := []struct {
		name            string
		method          string
		path            string
		booksErr        error
		usersErr        error
		wantStatus      int
		wantContentType string
		wantScript      string // какой модуль страницы подключён — по нему видно, что отдан нужный файл
	}{
		{name: "home", method: http.MethodGet, path: "/", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/home.js"},
		{name: "catalog", method: http.MethodGet, path: "/books?title=x", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/books.js"},
		{name: "book", method: http.MethodGet, path: "/books/3", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/book.js"},
		{name: "readers", method: http.MethodGet, path: "/users", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/users.js"},
		{name: "profile", method: http.MethodGet, path: "/users/4", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/user.js"},
		{name: "stats", method: http.MethodGet, path: "/stats", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/stats.js"},
		{name: "signup", method: http.MethodGet, path: "/signup", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/signup.js"},
		{name: "about", method: http.MethodGet, path: "/about", wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/static.js"},
		{name: "unknown page", method: http.MethodGet, path: "/nope", wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8", wantScript: "pages/static.js"},
		{name: "wrong method on home", method: http.MethodPost, path: "/", wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8"},
		// для путей API — прежний текстовый 404, а не HTML-страница
		{name: "unknown api path", method: http.MethodGet, path: "/api/v2/x", wantStatus: http.StatusNotFound, wantContentType: "text/plain; charset=utf-8"},
		{name: "script", method: http.MethodGet, path: "/static/js/ui.js", wantStatus: http.StatusOK, wantContentType: "text/javascript; charset=utf-8"},
		{name: "stylesheet", method: http.MethodGet, path: "/static/css/main.css", wantStatus: http.StatusOK, wantContentType: "text/css; charset=utf-8"},
		{name: "book not found", method: http.MethodGet, path: "/books/999999", booksErr: core_errors.ErrNotFound, wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8", wantScript: "pages/static.js"},
		{name: "book id not a number", method: http.MethodGet, path: "/books/abc", wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8"},
		{name: "book id zero", method: http.MethodGet, path: "/books/0", wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8"},
		{name: "book id above int4", method: http.MethodGet, path: "/books/3000000000", wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8"},
		// база недоступна — страница всё равно отдаётся, ошибку покажет JS
		{name: "book service failure", method: http.MethodGet, path: "/books/3", booksErr: errors.New("db is down"), wantStatus: http.StatusOK, wantContentType: "text/html; charset=utf-8", wantScript: "pages/book.js"},
		{name: "user not found", method: http.MethodGet, path: "/users/999999", usersErr: core_errors.ErrNotFound, wantStatus: http.StatusNotFound, wantContentType: "text/html; charset=utf-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &core_logger.Logger{Logger: zap.NewNop()}
			srv := core_http_server.NewHTTPServer(core_http_server.Config{}, log, core_http_middleware.Logger(log))
			srv.RegisterStatic("/static/", Static())
			srv.RegisterRoutes(Routes(fakeBooks{err: tt.booksErr}, fakeUsers{err: tt.usersErr})...)

			rec := httptest.NewRecorder()
			srv.Handler().ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Content-Type"); got != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", got, tt.wantContentType)
			}
			if tt.wantScript != "" && !strings.Contains(rec.Body.String(), tt.wantScript) {
				t.Errorf("body does not include %q", tt.wantScript)
			}
		})
	}
}
