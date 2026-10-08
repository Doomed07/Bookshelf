package users_transport_http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	core_auth "github.com/Doomed07/Bookshelf/internal/core/auth"
	core_domain "github.com/Doomed07/Bookshelf/internal/core/domain"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	"go.uber.org/zap"
)

// spy — «следующий обработчик»: запоминает, дошёл ли до него запрос.
type spy struct{ called bool }

func (s *spy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.called = true
	w.WriteHeader(http.StatusOK)
}

// newRequest кладёт в контекст логгер: в проде это делает глобальный middleware Logger.
func newRequest(method, target, body string) *http.Request {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	return r.WithContext(core_logger.ToCtx(r.Context(), &core_logger.Logger{Logger: zap.NewNop()}))
}

var pathValues = strings.NewReplacer("{id}", "7", "{user_id}", "7", "{book_id}", "3")

// Каждый маршрут, который что-то меняет, обязан быть закрыт: гость получает 401, чужой
// пользователь 403 и только владелец доходит до обработчика. Тест обходит ВСЕ маршруты,
// поэтому новый незащищённый PATCH/DELETE/POST провалит его сам.
func TestRoutesAreProtected(t *testing.T) {
	owner := core_auth.Identity{UserID: 7, Username: "owner", Email: "owner@mail.com"}
	stranger := core_auth.Identity{UserID: 8, Username: "stranger", Email: "stranger@mail.com"}

	routes := NewUsersHTTPHandler(nil).Routes()
	if len(routes) == 0 {
		t.Fatal("no routes")
	}

	for _, route := range routes {
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			next := &spy{}
			route.Handler = next.ServeHTTP // настоящий обработчик не нужен: проверяем только защиту
			mux := http.NewServeMux()
			mux.Handle(route.Method+" "+route.Path, route.WithMiddleware())
			target := pathValues.Replace(route.Path)

			serve := func(identity *core_auth.Identity) int {
				next.called = false
				r := newRequest(route.Method, target, "")
				if identity != nil {
					r = r.WithContext(core_auth.WithIdentity(r.Context(), *identity))
				}
				rec := httptest.NewRecorder()
				mux.ServeHTTP(rec, r)
				return rec.Code
			}

			if route.Method == http.MethodGet {
				if code := serve(nil); code != http.StatusOK || !next.called {
					t.Errorf("guest: status %d, handler called %v; reads must stay public", code, next.called)
				}
				return
			}

			if code := serve(nil); code != http.StatusUnauthorized || next.called {
				t.Errorf("guest: status %d, handler called %v; want 401 and not called", code, next.called)
			}
			if code := serve(&stranger); code != http.StatusForbidden || next.called {
				t.Errorf("another user: status %d, handler called %v; want 403 and not called", code, next.called)
			}
			if code := serve(&owner); code != http.StatusOK || !next.called {
				t.Errorf("owner: status %d, handler called %v; want 200 and called", code, next.called)
			}
		})
	}
}

// Регистрация переехала в /auth/register: публично создать пользователя без пароля здесь нельзя.
func TestNoPublicUserCreation(t *testing.T) {
	for _, route := range NewUsersHTTPHandler(nil).Routes() {
		if route.Method == http.MethodPost {
			t.Errorf("unexpected POST route %q: users must be created through /auth/register", route.Path)
		}
	}
}

type fakeUsersService struct {
	users []core_domain.User
	user  core_domain.User

	gotPatch core_domain.UserPatch
}

func (f *fakeUsersService) GetUsers(ctx context.Context, limit, offset *int) ([]core_domain.User, error) {
	return f.users, nil
}

func (f *fakeUsersService) GetUser(ctx context.Context, id int) (core_domain.User, error) {
	return f.user, nil
}

func (f *fakeUsersService) PatchUser(ctx context.Context, id int, patch core_domain.UserPatch) (core_domain.User, error) {
	f.gotPatch = patch
	return f.user, nil
}

func (f *fakeUsersService) DeleteUser(ctx context.Context, id int) error { return nil }

// Email — личные данные: ни один публичный ответ про пользователей не должен его содержать.
func TestResponsesDoNotExposeEmail(t *testing.T) {
	created := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	user := core_domain.NewUser(7, 3, "book_worm07", "secret@mail.com", created)
	svc := &fakeUsersService{user: user, users: []core_domain.User{user, core_domain.NewUser(8, 1, "anna", "anna@mail.com", created)}}
	h := NewUsersHTTPHandler(svc)

	tests := []struct {
		name    string
		handler http.HandlerFunc
		request *http.Request
		isList  bool
	}{
		{name: "GET /users", handler: h.GetUsers, request: newRequest(http.MethodGet, "/users", ""), isList: true},
		{name: "GET /users/{id}", handler: h.GetUser, request: pathID(newRequest(http.MethodGet, "/users/7", ""))},
		{name: "PATCH /users/{id}", handler: h.PatchUser, request: pathID(newRequest(http.MethodPatch, "/users/7", `{"email":"secret@mail.com"}`))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			tt.handler(rec, tt.request)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
			}
			body := rec.Body.String()
			for _, leak := range []string{"email", "secret@mail.com", "anna@mail.com", "@mail.com"} {
				if strings.Contains(body, leak) {
					t.Errorf("response leaks %q: %s", leak, body)
				}
			}

			// публичные поля остались на месте: id, version, username
			var objects []map[string]any
			if tt.isList {
				if err := json.Unmarshal([]byte(body), &objects); err != nil {
					t.Fatalf("decode list: %v; body = %s", err, body)
				}
			} else {
				var one map[string]any
				if err := json.Unmarshal([]byte(body), &one); err != nil {
					t.Fatalf("decode object: %v; body = %s", err, body)
				}
				objects = []map[string]any{one}
			}
			if len(objects) == 0 {
				t.Fatal("no users in the response")
			}
			for _, obj := range objects {
				if len(obj) != 3 || obj["id"] == nil || obj["version"] == nil || obj["username"] == nil {
					t.Errorf("user object = %v, want exactly id, version, username", obj)
				}
			}
		})
	}

	// владелец по-прежнему может сменить почту: запрос на неё доходит до сервиса
	if !svc.gotPatch.Email.Set || svc.gotPatch.Email.Value == nil || *svc.gotPatch.Email.Value != "secret@mail.com" {
		t.Errorf("service got patch %+v, want the email change to reach it", svc.gotPatch)
	}
}

func pathID(r *http.Request) *http.Request {
	r.SetPathValue("id", "7")
	return r
}
