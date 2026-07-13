package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/catstyle1101/todo_app_go/cmd/internal/core/config"
	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_pgx_pool "github.com/catstyle1101/todo_app_go/cmd/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/middleware"
	core_http_server "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/server"
	tasks_postgres_repository "github.com/catstyle1101/todo_app_go/cmd/internal/features/tasks/repository/postgres"
	tasks_service "github.com/catstyle1101/todo_app_go/cmd/internal/features/tasks/service"
	tasks_transport_http "github.com/catstyle1101/todo_app_go/cmd/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/catstyle1101/todo_app_go/cmd/internal/features/users/repository/postgres"
	users_service "github.com/catstyle1101/todo_app_go/cmd/internal/features/users/service"
	users_transport_http "github.com/catstyle1101/todo_app_go/cmd/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to load logger")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Application timezone", zap.Any("zone", time.Local))

	logger.Info("🚀 Starting ToDo application!")

	logger.Debug("initializing postgres pool...")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())

	if err != nil {
		logger.Fatal("error when initializing pool: %w", zap.Error(err))
	}
	logger.Debug("Postgres pool successfully initialized")
	defer pool.Close()

	logger.Debug("Initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	userService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHttpHandler(userService)

	logger.Debug("Initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTaskHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	/*
		Example of usage middleware on separate api version router

		apiVersionRouterV2 := core_http_server.NewApiVersionRouter(
			core_http_server.ApiVersion2,
			core_http_middleware.Dummy("example"),
		)
		apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)
	*/

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		// apiVersionRouterV2,
	)

	if err = httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server error: %w", zap.Error(err))
	}
}
