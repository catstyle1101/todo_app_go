package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("1")
	ApiVersion2 = ApiVersion("2")
	ApiVersion3 = ApiVersion("3")
)

type ApiVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleWare []core_http_middleware.MiddleWare
}

func NewApiVersionRouter(
	apiVersion ApiVersion,
	middleware ...core_http_middleware.MiddleWare,
) *ApiVersionRouter {
	return &ApiVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleWare: middleware,
	}
}

func (r *ApiVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *ApiVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(r, r.middleWare...)
}
