package store

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	client *redis.Client
}

func NewRedisDB(ctx context.Context, address string) (Store, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt ping redis")
	}

	return &redisDB{
		client: rdb,
	}, nil
}

func (e *redisDB) GetAllTasks(ctx context.Context) ([]Task, error) {
	keys, err := e.client.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, errors.Wrap(err, "couldnt query for tasks")
	}

	tasks := []Task{}

	for _, id := range keys {
		val, err := e.client.Get(ctx, id).Result()
		if err != nil {
			return nil, errors.Wrapf(err, "couldnt query for task %s", id)
		}

		t := Task{}
		err = json.Unmarshal([]byte(val), &t)
		if err != nil {
			return nil, errors.Wrap(err, "couldnt parsing tasks from string")
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (e *redisDB) CreateTask(ctx context.Context, task Task) (*Task, error) {
	task.ID = uuid.NewString()

	value, err := json.Marshal(task)
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt json marshal task %s", task.ID)
	}

	err = e.client.Set(ctx, "task:"+task.ID, string(value), 0).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt create task %s", task.ID)
	}

	return &task, nil
}
func (e *redisDB) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	if id == "" {
		return nil, errors.Errorf("there is no id provided")
	}

	val, err := e.client.Get(ctx, "task:"+id).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt query for task %s", id)
	}

	t := Task{}
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt parsing task from string")
	}

	return &t, nil
}

func (e *redisDB) UpdateTask(ctx context.Context, task Task) (*Task, error) {
	value, err := json.Marshal(task)
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt json marshal task %s", task.ID)
	}

	err = e.client.Set(ctx, "task:"+task.ID, string(value), 0).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt update task %s", task.ID)
	}

	return &task, nil
}

func (e *redisDB) DeleteTaskByID(ctx context.Context, id string) (*Task, error) {
	val, err := e.client.Get(ctx, "task:"+id).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt query for task %s", id)
	}

	t := Task{}
	err = json.Unmarshal([]byte(val), &t)
	if err != nil {
		return nil, errors.Wrap(err, "couldnt parsing task from string")
	}

	err = e.client.Del(ctx, "task:"+t.ID).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "couldnt delete task %s", t.ID)
	}

	return &t, nil
}
