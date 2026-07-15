package tasks_transport_http

import (
	"time"

	"github.com/catstyle1101/todo_app_go/cmd/internal/core/domain"
)

type TaskDtoResponse struct {
	ID           int        `json:"id"             example:"333"`
	Version      int        `json:"version"        example:"4"`
	Title        string     `json:"title"          example:"Домашка"`
	Description  *string    `json:"description"    example:"Сделать до четверга домашнее задание по математике"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-07-30T10:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at"   example:"null"`
	AuthorUserID int        `json:"author_user_id" example:"5"`
}

func taskDTOFromDomain(task domain.Task) TaskDtoResponse {
	return TaskDtoResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDtoResponse {
	dtos := make([]TaskDtoResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
