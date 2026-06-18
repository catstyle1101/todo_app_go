package users_transport_http

import (
	"net/http"

	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
	core_http_utils "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := core_http_utils.GetLimitOffsetQueryParams(r)

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
