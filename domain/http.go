package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/EarvinKayonga/tasks/store"
)

func ServeAPI(svc store.Store) func() error {
	return func() error {
		// getTaskByID returns the task with the correct ID.
		getTaskByID := func(w http.ResponseWriter, r *http.Request) {
			taskID := chi.URLParam(r, "taskID")

			task, err := svc.GetTaskByID(r.Context(), taskID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		// updateTaskByID update the task with the given ID
		// returns the updated task.
		updateTaskByID := func(w http.ResponseWriter, r *http.Request) {
			taskID := chi.URLParam(r, "taskID")

			task := &store.Task{}
			err := json.NewDecoder(r.Body).Decode(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)

				return
			}

			task.ID = taskID
			task, err = svc.UpdateTask(r.Context(), *task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		deleteTasksByID := func(w http.ResponseWriter, r *http.Request) {
			taskID := chi.URLParam(r, "taskID")

			task, err := svc.DeleteTaskByID(r.Context(), taskID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		getTasks := func(w http.ResponseWriter, r *http.Request) {
			tasks, err := svc.GetAllTasks(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			err = json.NewEncoder(w).Encode(tasks)
			if err != nil {
				fmt.Fprintf(w, "%v", err.Error())
			}
		}

		createTask := func(w http.ResponseWriter, r *http.Request) {
			task := &store.Task{}
			err := json.NewDecoder(r.Body).Decode(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)

				return
			}

			task, err = svc.CreateTask(r.Context(), *task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusCreated)
			err = json.NewEncoder(w).Encode(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
		}

		router := chi.NewRouter()

		router.Route("/tasks", func(r chi.Router) {
			r.Post("/", createTask)
			r.Get("/", getTasks)
			r.Get("/{taskID}", getTaskByID)
			r.Put("/{taskID}", updateTaskByID)
			r.Delete("/{taskID}", deleteTasksByID)
		})

		address := ":4000" // Vous pouvez aussi utiliser flag ou cli pour permettre de configurer l'adresse

		log.Printf("Listening on %s", address)
		err := http.ListenAndServe(address, router)
		if err != nil {
			return err
		}

		return nil
	}
}
