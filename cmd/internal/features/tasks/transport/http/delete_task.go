package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
)

func (h *TaskHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetPathIntValue(r, "task_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "get task_id path value")
		return
	}

	err = h.taskService.DeleteTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(err, "delete task")
		return
	}

	responseHandler.NoContentResponse()

}
