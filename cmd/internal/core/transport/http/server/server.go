package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_middleware "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/middleware"
	"github.com/catstyle1101/todo_app_go/docs"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.MiddleWare
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.MiddleWare,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)

		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := fmt.Sprintf("/api/v%s", router.apiVersion)
		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()),
		)
	}
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		s.mux.Handle(pattern, route.WithMiddleware())
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}
