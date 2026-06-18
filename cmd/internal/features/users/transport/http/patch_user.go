package users_transport_http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
	core_errors "github.com/catstyle1101/todo_app_go/cmd/internal/core/errors"
	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
	core_http_types "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/types"
	core_http_utils "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'full_name' can't be null")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("'full_name' must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf(
					"invalid `phone_number` len: %d, %w",
					phoneNumberLen,
					core_errors.ErrInvalidArgument,
				)
			}

			re := regexp.MustCompile(`^\+[0-9]{9,14}$`)

			if !re.MatchString(*r.PhoneNumber.Value) {
				return fmt.Errorf(
					"invalid `phone_number` format: %w",
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetPathIntValue(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user id path value")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	patchedUser, err := h.usersService.PatchUser(ctx, userID, userPatch)

	if err != nil {
		responseHandler.ErrorResponse(err, "patch user with service")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(patchedUser))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
