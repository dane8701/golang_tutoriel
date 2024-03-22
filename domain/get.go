package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/EarvinKayonga/tasks/store"
)

func ListTasks(svc store.Store) func(context.Context) error {
	return func(ctx context.Context) error {
		tasks, err := svc.GetAllTasks(ctx)
		if err != nil {
			return errors.Wrap(err, "couldnt get all tasks")
		}

		PrintTasks(tasks...)

		return nil
	}
}

func GetTaskByID(svc store.Store) func(context.Context, string) error {
	return func(ctx context.Context, taskID string) error {
		task, err := svc.GetTaskByID(ctx, taskID)
		if err != nil {
			return errors.Wrapf(err, "couldnt get task with %s", taskID)
		}

		PrintTasks(*task)

		return nil
	}
}

func PrintTasks(tasks ...store.Task) {
	numberOfTasks := len(tasks)

	if numberOfTasks == 0 {
		fmt.Println("no tasks to print")
		return
	}

	fmt.Printf("printing %d tasks", numberOfTasks)
	for _, task := range tasks {
		fmt.Println(jsonify(task))
	}
}

func jsonify(data store.Task) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
