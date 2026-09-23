package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Doomed07/Bookshelf/internal/core/logger"
	core_postgres_pgx "github.com/Doomed07/Bookshelf/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/Doomed07/Bookshelf/internal/core/transport/http/middleware"
	core_http_server "github.com/Doomed07/Bookshelf/internal/core/transport/http/server"
	users_repository_postgres "github.com/Doomed07/Bookshelf/internal/featurs/users/repository/postgres"
	users_service "github.com/Doomed07/Bookshelf/internal/featurs/users/service"
	users_transport_http "github.com/Doomed07/Bookshelf/internal/featurs/users/transport/http"
	"go.uber.org/zap"
)

func main() {
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

	logger.Debug("initializing connection pool")
	pool, err := core_postgres_pgx.NewConnPool(
		core_postgres_pgx.NewConfigMust(),
		ctx,
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "Users"))

	usersRepository := users_repository_postgres.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHttp := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHttp.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("Failed to run server", zap.Error(err))
	}
}
