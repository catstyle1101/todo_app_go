package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
	core_http_types "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/types"
)

type PatchTaskResponse TaskDtoResponse

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Погулять с собакой"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` cant be null")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))

			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be null")
		}
	}

	return nil
}

// PatchTask 	godoc
// @Summary 	Изменить задачу
// @Description Изменение информации об уже существующей в системе задачи
// @Description ### Логика обновления полей (three-state logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `{"description": "Выйти прогуляться"}` - устанавливает новое описание в БД
// @Description 3. **Передан `null`**: `{"description": null}` - очищает поле в БД (set NULL)
// @Description Ограничения: `title` и `completed` не могут быть выставлены как `null`.
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		task_id			path		int		true 					"ID задачи для обновления"
// @Param		request body	PatchTaskRequest 	true 			"PatchTaskRequest тело запроса"
// @Success 	200 {object} 	PatchTaskResponse 					"Успешно измененная задача"
// @Failure 	400 {object} 	core_http_response.ErrorResponse 	"Bad Request"
// @Failure 	404 {object} 	core_http_response.ErrorResponse 	"User not found"
// @Failure 	409 {object} 	core_http_response.ErrorResponse 	"Conflict"
// @Failure 	500 {object} 	core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/tasks/{task_id} [patch]
func (h *TaskHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetPathIntValue(r, "task_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "get task_id path value")
		return
	}

	var request PatchTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "decode and validate HTTP request")
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.taskService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
