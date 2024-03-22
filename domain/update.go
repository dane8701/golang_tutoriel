package domain

import (
	"context"

	"github.com/EarvinKayonga/tasks/store"
)

func UpdateTaskByID(svc store.Store) func(context.Context, string, string, store.Status) error {
	return func(ctx context.Context, id, name string, status store.Status) error {
		task, err := svc.UpdateTask(ctx, store.Task{
			ID:     id,
			Status: status,
			Name:   name,
		})
		if err != nil {
			return err
		}

		PrintTasks(*task)

		return nil
	}
}
