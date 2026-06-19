package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/catstyle1101/todo_app_go/cmd/internal/core/errors"
)

var (
	LimitQueryParamKey  = "limit"
	OffsetQueryParamKey = "offset"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)

	if err != nil {
		return nil, fmt.Errorf(
			"param='%s', by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := GetIntQueryParam(r, LimitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := GetIntQueryParam(r, OffsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
