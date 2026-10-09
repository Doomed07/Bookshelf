package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_auth "github.com/Doomed07/Bookshelf/internal/core/auth"
	core_config "github.com/Doomed07/Bookshelf/internal/core/config"
	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_metrics "github.com/Doomed07/Bookshelf/internal/core/metrics"
	core_postgres_pgx "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Doomed07/Bookshelf/internal/core/transport/http/middleware"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
	auth_repository_postgres "github.com/Doomed07/Bookshelf/internal/features/auth/repository/postgres"
	auth_service "github.com/Doomed07/Bookshelf/internal/features/auth/service"
	auth_transport_http "github.com/Doomed07/Bookshelf/internal/features/auth/transport/http"
	books_repository_postgres "github.com/Doomed07/Bookshelf/internal/features/books/repository/postgres"
	books_service "github.com/Doomed07/Bookshelf/internal/features/books/service"
	books_transport_http "github.com/Doomed07/Bookshelf/internal/features/books/transport/http"
	bookshelf_repository_postgres "github.com/Doomed07/Bookshelf/internal/features/bookshelf/repository/postgres"
	bookshelf_service "github.com/Doomed07/Bookshelf/internal/features/bookshelf/service"
	bookshelf_transport_http "github.com/Doomed07/Bookshelf/internal/features/bookshelf/transport/http"
	statistics_repository_postgres "github.com/Doomed07/Bookshelf/internal/features/statistics/repository/postgres"
	statistics_service "github.com/Doomed07/Bookshelf/internal/features/statistics/service"
	statistics_transport_http "github.com/Doomed07/Bookshelf/internal/features/statistics/transport/http"
	users_repository_postgres "github.com/Doomed07/Bookshelf/internal/features/users/repository/postgres"
	users_service "github.com/Doomed07/Bookshelf/internal/features/users/service"
	users_transport_http "github.com/Doomed07/Bookshelf/internal/features/users/transport/http"
	"github.com/Doomed07/Bookshelf/web"
	"go.uber.org/zap"

	_ "github.com/Doomed07/Bookshelf/docs"
)

// @title 			Bookshelf API
// @version			1.0
// @description 		Bookshelf Application REST-API schema
// @host 			127.0.0.1:8080
// @BasePath 		/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone
	statsConfig := statistics_service.NewConfigMust()
	authConfig := core_auth.NewConfigMust()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing connection pool")
	pool, err := core_postgres_pgx.NewConnPool(
		core_postgres_pgx.NewConfigMust(),
		ctx,
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()
	core_metrics.RegisterPool(pool.Pool)

	logger.Debug("initializing feature", zap.String("feature", "Users"))
	usersRepository := users_repository_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "Books"))
	booksRepository := books_repository_postgres.NewBooksRepository(pool)
	booksService := books_service.NewBooksService(booksRepository)
	booksTransportHTTP := books_transport_http.NewBooksHTTPHandler(booksService)

	logger.Debug("initializing feature", zap.String("feature", "Bookshelf"))
	bookshelfRepository := bookshelf_repository_postgres.NewBookshelfRepository(pool)
	bookshelfService := bookshelf_service.NewBookshelfService(bookshelfRepository)
	bookshelfTransportHTTP := bookshelf_transport_http.NewBookshelfHTTPHandler(bookshelfService)

	logger.Debug("initializing feature", zap.String("feature", "Statistics"))
	statsRepository := statistics_repository_postgres.NewStatsRepository(pool)
	statsService := statistics_service.NewStatsService(statsRepository, statsConfig)
	statsTransportHTTP := statistics_transport_http.NewStatsHTTPHandler(statsService)

	logger.Debug("initializing feature", zap.String("feature", "Auth"))
	authRepository := auth_repository_postgres.NewAuthRepository(pool)
	authService, err := auth_service.NewAuthService(authRepository, authConfig.BcryptCost, authConfig.SessionTTL)
	if err != nil {
		logger.Fatal("failed to init auth service", zap.Error(err))
	}
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService, authConfig.SessionTTL, authConfig.CookieSecure)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.Metrics(),
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(
		core_http_server.ApiVersion1,
		core_http_middleware.SameOrigin(),
		core_http_middleware.Authenticate(authService),
	)

	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(booksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(bookshelfTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statsTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(authTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)

	httpServer.RegisterSwagger()

	httpServer.RegisterStatic("/static/", web.Static())
	httpServer.RegisterRoutes(web.Routes(booksService, usersService)...)

	logger.Debug("initializing metrics server")
	metricsServer := core_http_server.NewHTTPServer(
		core_http_server.NewMetricsConfigMust(),
		logger,
	)
	metricsServer.RegisterRoutes(core_http_server.Route{
		Method:  http.MethodGet,
		Path:    "/metrics",
		Handler: core_metrics.Handler().ServeHTTP,
	})

	go func() {
		if err := metricsServer.Run(ctx); err != nil {
			logger.Error("Failed to run metrics server", zap.Error(err))
		}
	}()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Failed to run server", zap.Error(err))
	}
}
