package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/koliashka/golang-todoapp/internal/core/logger"
	core_postgres_pool "github.com/koliashka/golang-todoapp/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/koliashka/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/koliashka/golang-todoapp/internal/core/transport/server"
	users_postgres_repository "github.com/koliashka/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/koliashka/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/koliashka/golang-todoapp/internal/features/users/transport/http"
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
		fmt.Println("failed to init app loigger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres conn pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initioalizing feature", zap.String("feature", "users"))
	userRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(userRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("init HTTP serveice")

	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error: %w", zap.Error(err))
	}
}
