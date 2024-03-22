package domain

import (
	"context"

	"github.com/EarvinKayonga/tasks/store"
)

func CreateTask(svc store.Store) func(context.Context, string, store.Status) error {
	return func(ctx context.Context, name string, status store.Status) error {
		task, err := svc.CreateTask(ctx, store.Task{
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
