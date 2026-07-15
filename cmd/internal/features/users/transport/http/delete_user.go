package users_transport_http

import (
	"net/http"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
)

// DeleteUser 	godoc
// @Summary 	Удалить пользователя
// @Description Удалить пользователя из системы по его ID
// @Tags 		users
// @Param		user_id path int true "ID удаляемого пользователя"
// @Success 	204 "Успешное удаление пользователя"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 	404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router 		/users/{user_id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetPathIntValue(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user id path value")
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
