package statistics_postgres_repository

import core_postgres_pool "github.com/catstyle1101/todo_app_go/cmd/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatisticsRepository(pool core_postgres_pool.Pool) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
