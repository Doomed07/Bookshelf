package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
)

func Test(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromCtx(ctx)

			log.Debug(fmt.Sprintf("-> input s: %s", s))

			next.ServeHTTP(w, r)

			log.Debug(fmt.Sprintf("<- output s: %s", s))
		})
	}
}
