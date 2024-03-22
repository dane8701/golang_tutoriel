package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type jsonFile struct {
	file string
}

func NewJsonFile(file string) (Store, error) {
	info, err := os.Stat(file)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt stat file")
	}

	if info.IsDir() {
		return nil, fmt.Errorf("%s is a folder", file)
	}

	return &jsonFile{file}, nil
}

func (e *jsonFile) CreateTask(ctx context.Context, task Task) (*Task, error) {
	tasks, err := e.loadTasks()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt load tasks from file")
	}

	task.ID = uuid.New().String()
	tasks = append(tasks, task)

	err = e.writeTasks(tasks)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt load tasks from file")
	}

	return &task, nil
}

func (e *jsonFile) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("empty id provided")
	}

	tasks, err := e.loadTasks()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt load tasks from file")
	}

	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}

	return nil, fmt.Errorf("task %s not found", id)
}
func (e *jsonFile) GetAllTasks(ctx context.Context) ([]Task, error) {
	tasks, err := e.loadTasks()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt load tasks from file")
	}

	return tasks, nil
}
func (e *jsonFile) UpdateTask(ctx context.Context, task Task) (*Task, error) {
	tasks, err := e.loadTasks()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt load tasks from file")
	}

	for i := range tasks {
		if tasks[i].ID == task.ID {
			tasks[i] = task
			err = e.writeTasks(tasks)
			if err != nil {
				return nil, errors.Wrap(err, "couldnt load tasks from file")
			}

			return &task, nil
		}
	}

	return nil, fmt.Errorf("task %s not found", task.ID)
}

func (e *jsonFile) DeleteTaskByID(ctx context.Context, id string) (*Task, error) {
	tasks, err := e.loadTasks()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt load tasks from file")
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			err = e.writeTasks(tasks)
			if err != nil {
				return nil, errors.Wrap(err, "couldnt load tasks from file")
			}

			return &tasks[i], nil
		}
	}

	return nil, fmt.Errorf("task %s not found", id)
}

func (e *jsonFile) loadTasks() ([]Task, error) {
	file, err := os.Open(e.file)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt open file")
	}

	defer file.Close()

	tasks := []Task{}
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt decode task from file")
	}

	return tasks, nil
}

func (e *jsonFile) writeTasks(tasks []Task) error {
	file, err := os.Create(e.file)
	if err != nil {
		return errors.Wrap(err, "couldnt open file")
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return errors.Wrap(err, "couldnt decode task from file")
	}

	return nil
}
