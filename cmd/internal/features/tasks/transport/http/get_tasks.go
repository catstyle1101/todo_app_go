package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDtoResponse

func (h *TaskHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	user_id, err := core_http_request.GetIntQueryParam(r, "user_id")

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'user_id' query param")
		return
	}

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit'/'offset' query param")
		return
	}

	tasksDomains, err := h.taskService.GetTasks(ctx, user_id, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "get tasks")
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(response, http.StatusOK)

}
