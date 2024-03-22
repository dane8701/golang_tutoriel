package domain

import (
	"context"

	"github.com/pkg/errors"

	"github.com/EarvinKayonga/tasks/store"
)

func DeleteTaskByID(svc store.Store) func(context.Context, string) error {
	return func(ctx context.Context, taskID string) error {
		task, err := svc.DeleteTaskByID(ctx, taskID)
		if err != nil {
			return errors.Wrapf(err, "couldnt delete task with %s", taskID)
		}

		PrintTasks(*task)

		return nil
	}
}
