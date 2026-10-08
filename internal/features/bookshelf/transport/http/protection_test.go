package bookshelf_transport_http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	core_auth "github.com/Doomed07/Bookshelf/internal/core/auth"
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
func newRequest(method, target string) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	return r.WithContext(core_logger.ToCtx(r.Context(), &core_logger.Logger{Logger: zap.NewNop()}))
}

var pathValues = strings.NewReplacer("{id}", "7", "{user_id}", "7", "{book_id}", "3")

// Полку может менять только её владелец. Тест обходит ВСЕ маршруты: гость получает 401,
// другой пользователь 403, владелец доходит до обработчика. Чтение полки остаётся открытым.
// Новый незащищённый POST/PATCH/DELETE провалит этот тест сам.
func TestRoutesAreProtected(t *testing.T) {
	owner := core_auth.Identity{UserID: 7, Username: "owner"}
	stranger := core_auth.Identity{UserID: 8, Username: "stranger"}

	routes := NewBookshelfHTTPHandler(nil).Routes()
	if len(routes) == 0 {
		t.Fatal("no routes")
	}

	writes := 0
	for _, route := range routes {
		t.Run(route.Method+" "+route.Path, func(t *testing.T) {
			next := &spy{}
			route.Handler = next.ServeHTTP // настоящий обработчик не нужен: проверяем только защиту
			mux := http.NewServeMux()
			mux.Handle(route.Method+" "+route.Path, route.WithMiddleware())
			target := pathValues.Replace(route.Path)

			serve := func(identity *core_auth.Identity) int {
				next.called = false
				r := newRequest(route.Method, target)
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
		if route.Method != http.MethodGet {
			writes++
		}
	}

	// страховка от тихой потери защиты: добавление, изменение и удаление книги на полке
	if writes != 3 {
		t.Errorf("write routes = %d, want 3 (add, patch, remove)", writes)
	}
}
