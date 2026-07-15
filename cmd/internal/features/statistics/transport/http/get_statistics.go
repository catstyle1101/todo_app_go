package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
	core_logger "github.com/catstyle1101/todo_app_go/cmd/internal/core/logger"
	core_http_request "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/request"
	core_http_response "github.com/catstyle1101/todo_app_go/cmd/internal/core/transport/http/response"
)

type StatisticsResponse struct {
	TasksCreated              int      `json:"tasks_created" example:"30"`
	TasksCompleted            int      `json:"tasks_completed" example:"10"`
	TasksCompletedRate        *float64 `json:"tasks_completed_rate" example:"33.33333"`
	TasksAverageCompletedTime *string  `json:"tasks_average_completion_time" example:"5h32m12s"`
}

// GetStatistics godoc
// @Summary 	Получение статистики
// @Summary 	Получение статистики по задачам с опциональной фильтрацией по user_id и/или по временному промежутку создания задач
// @Tags 		statistics
// @Produce 	json
// @Param 		user_id 	query 	int 	false 	"Фильтрация статистики по конкретному пользователю"
// @Param 		from 		query 	int 	false 	"Начало промежутка рассмотрения статистики (включительно), формат `YYYY-MM-DD`"
// @Param 		to 			query 	int 	false 	"Конец промежутка рассмотрения статистики (не включительно), формат `YYYY-MM-DD`"
// @Success 	200 		{object} 		StatisticsResponse 					"Успешное получение статистики"
// @Failure 	400 		{object} 		core_http_response.ErrorResponse 	"Bad request"
// @Failure 	500 		{object} 		core_http_response.ErrorResponse 	"Internal server error"
// @Router 		/statistics 				[get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := GetUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/from/to query params")
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := toDtoFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func toDtoFromDomain(statistics domain.Statistics) StatisticsResponse {
	var avgCompletedTime *string
	if statistics.TasksAverageCompletedTime != nil {
		duration := statistics.TasksAverageCompletedTime.String()
		avgCompletedTime = &duration
	}

	return StatisticsResponse{
		TasksCreated:              statistics.TasksCreated,
		TasksCompleted:            statistics.TasksCompleted,
		TasksCompletedRate:        statistics.TasksCompletedRate,
		TasksAverageCompletedTime: avgCompletedTime,
	}
}

func GetUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `user_id` query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `from` query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get `to` query param: %w", err)
	}

	return userID, from, to, nil

}
