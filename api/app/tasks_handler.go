package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// TasksHandlerGetAll returns all tasks for the current user
func TasksHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		j, err := json.Marshal(user.Tasks)
		if err != nil {
			mlog.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
	}
}

// TasksHandlerCreate creates a task for the current user, followed by a defrag of the user's tasks
// Requires a JSON body with the task's name, description, due, priority, and completed status
func TasksHandlerCreate(w http.ResponseWriter, r *http.Request) {
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		decoder := json.NewDecoder(r.Body)
		var t struct {
			Name        string
			Description string
			Due         time.Time
			Priority    int
			Completed   bool
		}
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if t.Name != "" && t.Description != "" && t.Priority > 0 {
			err = CreateTask(user, t.Name, t.Description, t.Due, t.Priority, t.Completed)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = DefragTasks(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprint(w, "Successfully created task")
		}
	}
}

// TasksHandlerGet returns a single requested task for the current user
func TasksHandlerGet(w http.ResponseWriter, r *http.Request) {
	taskid, err := strconv.Atoi(mux.Vars(r)["taskid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		for _, t := range user.Tasks {
			if t.Priority == taskid {
				j, _ := json.Marshal(t)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(j))
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}

// TasksHandlerUpdate handles updating a task's information, as well as reordering tasks if a priority changes
func TasksHandlerUpdate(w http.ResponseWriter, r *http.Request) {
	taskid, err := strconv.Atoi(mux.Vars(r)["taskid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		decoder := json.NewDecoder(r.Body)
		var t struct {
			Name        string
			Description string
			Due         time.Time
			Priority    int
			Completed   bool
		}
		err := decoder.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, task := range user.Tasks {
			if task.Priority == taskid {
				if t.Name != "" && t.Name != task.Name {
					task.UpdateName(t.Name)
				}
				if t.Description != "" && t.Description != task.Description {
					task.UpdateDescription(t.Name)
				}
				if !t.Due.IsZero() && t.Due != task.DueDate {
					task.UpdateDue(t.Due)
				}
				if t.Completed != task.Completed {
					task.UpdateCompleted(t.Completed)
				}
				// Priority must be updated last, as it could affect the other updates
				if t.Priority > 0 && t.Priority != task.Priority {
					task.UpdatePriority(user, t.Priority)
				}
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}

// TasksHandlerDelete deletes the specified task, all associated reports, and then runs a defrag on the task list
func TasksHandlerDelete(w http.ResponseWriter, r *http.Request) {
	taskid, err := strconv.Atoi(mux.Vars(r)["taskid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		for _, t := range user.Tasks {
			if t.Priority == taskid {
				err = t.Delete()
				if err != nil {
					http.Error(w, "400 Bad Request", http.StatusBadRequest)
					return
				}
				err = DefragTasks(user)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}
