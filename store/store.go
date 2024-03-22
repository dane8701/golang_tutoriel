package store

import (
	"context"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Status    Status    `json:"status"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Status string

var (
	StatusCompleted Status = "completed"
	StatusWIP       Status = "wip"
	StatusTODO      Status = "todo"
)

type Store interface {
	CreateTask(ctx context.Context, task Task) (*Task, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	UpdateTask(ctx context.Context, task Task) (*Task, error)
	DeleteTaskByID(ctx context.Context, id string) (*Task, error)
}
