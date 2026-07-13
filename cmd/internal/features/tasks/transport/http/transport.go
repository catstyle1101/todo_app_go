package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
	core_http_server "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/server"
)

type TaskHTTPHandler struct {
	taskService TaskService
}

type TaskService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
	PatchTask(ctx context.Context, id int, patch domain.TaskPatch) (domain.Task, error)
}

func NewTaskHTTPHandler(
	service TaskService,
) *TaskHTTPHandler {
	return &TaskHTTPHandler{
		taskService: service,
	}
}

func (h *TaskHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.GetTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{task_id}",
			Handler: h.GetTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{task_id}",
			Handler: h.DeleteTask,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{task_id}",
			Handler: h.PatchTask,
		},
	}
}
