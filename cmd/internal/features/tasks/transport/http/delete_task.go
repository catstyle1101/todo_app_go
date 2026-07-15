package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
)

// DeleteTask 	godoc
// @Summary 	Удалить задачу
// @Description Удалить задачу из системы по его ID
// @Tags 		tasks
// @Param		task_id path int true "ID удаляемой задачи"
// @Success 	204 "Успешное удаление пользователя"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router 		/tasks/{task_id} [delete]
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
