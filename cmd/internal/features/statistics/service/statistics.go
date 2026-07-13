package statistics_service

import (
	"context"
	"time"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(repository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: repository,
	}
}
