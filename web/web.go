// Package web хранит фронтенд Shelfmate: HTML-страницы, CSS, картинки и JS.
// Страницы — статичные каркасы, данные в них подгружает JS из /api/v1.
// Файлы встраиваются в бинарник, поэтому в Docker-образ их копировать не нужно.
package web

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"math"
	"net/http"
	"strings"

	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_errors "github.com/Doomed07/Bookshelf/internal/core/errors"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_http_request "github.com/Doomed07/Bookshelf/internal/core/transport/http/request"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
	"go.uber.org/zap"
)

//go:embed static
var files embed.FS

// Static — всё содержимое web/static, раздаётся по /static/.
func Static() fs.FS {
	static, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	return static
}

type BooksService interface {
	GetBook(ctx context.Context, id int) (core_domain.Book, error)
}

type UsersService interface {
	GetUser(ctx context.Context, id int) (core_domain.User, error)
}

// existing отдаёт страницу, только если запись с id из адреса существует.
// Иначе — 404.html со статусом 404, чтобы поисковики не считали страницу живой.
func existing(lookup func(ctx context.Context, id int) error, name string) http.HandlerFunc {
	ok := page(name, http.StatusOK)
	notFoundPage := page("404.html", http.StatusNotFound)

	return func(rw http.ResponseWriter, r *http.Request) {
		id, err := core_http_request.GetIntPathParam(r, "id")
		if err != nil || id < 1 || id > math.MaxInt32 { // больше int4 Postgres не примет
			notFoundPage(rw, r)
			return
		}

		if err := lookup(r.Context(), id); err != nil {
			if errors.Is(err, core_errors.ErrNotFound) {
				notFoundPage(rw, r)
				return
			}
			// база недоступна и т.п. — страницу всё равно отдаём, JS покажет ошибку
			core_logger.FromCtx(r.Context()).Error("check page entity", zap.Error(err))
		}

		ok(rw, r)
	}
}

// Routes — адреса страниц сайта. Каждая страница — HTML-файл; id из адреса
// вроде /books/3 JS читает сам из location.pathname.
func Routes(books BooksService, users UsersService) []core_http_server.Route {
	return []core_http_server.Route{
		{
			// {$} — только сам "/": шаблон "GET /" конфликтовал бы с "/api/v1/"
			Method:  http.MethodGet,
			Path:    "/{$}",
			Handler: page("index.html", http.StatusOK),
		},
		{
			Method:  http.MethodGet,
			Path:    "/books",
			Handler: page("books.html", http.StatusOK),
		},
		{
			Method: http.MethodGet,
			Path:   "/books/{id}",
			Handler: existing(func(ctx context.Context, id int) error {
				_, err := books.GetBook(ctx, id)
				return err
			}, "book.html"),
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: page("users.html", http.StatusOK),
		},
		{
			Method: http.MethodGet,
			Path:   "/users/{id}",
			Handler: existing(func(ctx context.Context, id int) error {
				_, err := users.GetUser(ctx, id)
				return err
			}, "user.html"),
		},
		{
			Method:  http.MethodGet,
			Path:    "/stats",
			Handler: page("stats.html", http.StatusOK),
		},
		{
			Method:  http.MethodGet,
			Path:    "/signup",
			Handler: page("signup.html", http.StatusOK),
		},
		{
			Method:  http.MethodGet,
			Path:    "/login",
			Handler: page("login.html", http.StatusOK),
		},
		{
			Method:  http.MethodGet,
			Path:    "/about",
			Handler: page("about.html", http.StatusOK),
		},
		{
			// без метода: всё, что не совпало ни с чем другим
			Path:    "/",
			Handler: notFound(page("404.html", http.StatusNotFound)),
		},
	}
}

func page(name string, statusCode int) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := fs.ReadFile(files, "static/"+name)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.WriteHeader(statusCode)
		_, _ = rw.Write(body)
	}
}

// notFound: неизвестные пути API отвечают обычным текстовым 404, как и раньше.
func notFound(htmlPage http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api" || strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(rw, r)
			return
		}

		htmlPage(rw, r)
	}
}
