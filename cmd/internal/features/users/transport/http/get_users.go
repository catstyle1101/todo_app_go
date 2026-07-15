package users_transport_http

import (
	"net/http"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers 	godoc
// @Summary 	Список пользователей
// @Description Получение списка пользователей с опциональной пагинацией
// @Tags 		users
// @Produce 	json
// @Param 		limit 	query 	int 	false				 	"Размер страницы с пользователями"
// @Param 		offset 	query 	int 	false 					"Смещение страницы с пользователями"
// @Success 	200 	{object} 		GetUsersResponse 					"Полученный список пользователей"
// @Failure 	400 	{object} 		core_http_response.ErrorResponse 	"Bad request"
// @Failure 	500 	{object} 		core_http_response.ErrorResponse 	"Internal Server Error"
// @Router 		/users 	[get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)

	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit'/'offset' query param")
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomain(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}
