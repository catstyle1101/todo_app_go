package tasks_service

import (
	"context"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID *int, limit *int, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, id int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
	PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error)
}

func NewTasksService(repository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: repository,
	}
}
