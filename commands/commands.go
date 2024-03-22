package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	"github.com/EarvinKayonga/tasks/domain"
	"github.com/EarvinKayonga/tasks/store"
)

func Create() *cli.App {
	return &cli.App{
		Name: "tasks",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "file",
				Value:    "",
				Usage:    "local tasks json file: try 'tmp/tasks.json'",
				Required: false,
			},
			cli.StringFlag{
				Name:     "redis",
				Value:    "",
				Usage:    "address to redis",
				Required: false,
			},
		},
		Commands: []cli.Command{
			{
				Name:  "serve",
				Usage: "Start the API server",
				Action: func(c *cli.Context) error {
					svc, err := store.NewRedisDB(context.Background(), c.GlobalString("redis"))
					if err != nil {
						return err
					}

					return domain.ServeAPI(svc)()
				},
			},
			{
				Name:  "list",
				Usage: "list all tasks",
				Action: func(c *cli.Context) error {
					svc, err := store.NewRedisDB(context.Background(), c.GlobalString("redis"))
					if err != nil {
						return err
					}

					return domain.ListTasks(svc)(context.Background())
				},
			},
			{
				Name:  "get",
				Usage: "get a task by its ID",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "id",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					svc, err := store.NewRedisDB(context.Background(), c.GlobalString("redis"))
					if err != nil {
						return err
					}

					taskID := c.String("id")

					if len(taskID) < 1 {
						fmt.Println("Please provide a task ID")
						return nil
					}

					return domain.GetTaskByID(svc)(context.Background(), taskID)
				},
			},
			{
				Name:  "create",
				Usage: "create a new task",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "name",
						Value:    "",
						Required: true,
					},
					cli.StringFlag{
						Name:     "status",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					svc, err := store.NewRedisDB(context.Background(), c.GlobalString("redis"))
					if err != nil {
						return err
					}

					name := c.String("name")
					status := c.String("status")
					if name == "" {
						fmt.Println("Please provide a task name")
						return nil
					}
					if status == "" {
						status = "todo" // Statut par défaut si non spécifié
					}
					return domain.CreateTask(svc)(context.Background(), name, store.Status(status))
				},
			},
			{
				Name:  "update",
				Usage: "update an existing task by its ID",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "name",
						Value:    "",
						Required: true,
					},
					cli.StringFlag{
						Name:     "status",
						Value:    "",
						Required: true,
					},
					cli.StringFlag{
						Name:     "id",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					svc, err := store.NewRedisDB(context.Background(), c.GlobalString("redis"))
					if err != nil {
						return err
					}

					taskID := c.String("id")
					newName := c.String("name")
					newStatus := store.Status(c.String("status"))
					return domain.UpdateTaskByID(svc)(context.Background(), taskID, newName, newStatus)
				},
			},
			{
				Name:  "delete",
				Usage: "delete an existing task by its ID",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     "id",
						Value:    "",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					svc, err := store.NewRedisDB(context.Background(), c.GlobalString("redis"))
					if err != nil {
						return err
					}

					taskID := c.String("id")

					if len(taskID) < 1 {
						fmt.Println("Please provide a task ID")
						return nil
					}
					return domain.DeleteTaskByID(svc)(context.Background(), taskID)
				},
			},
		},
	}
}
