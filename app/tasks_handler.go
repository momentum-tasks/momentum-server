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
		fmt.Fprint(w, string(j))
	}
}

// TasksHandlerCreate creates a task for the current user
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
		if t.Name != "" && t.Description != "" && t.Priority != 0 {
			err = CreateTask(user, t.Name, t.Description, t.Due, t.Priority, t.Completed)
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
				fmt.Fprint(w, string(j))
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}

func TasksHandlerUpdate(w http.ResponseWriter, r *http.Request) {
}

func TasksHandlerDelete(w http.ResponseWriter, r *http.Request) {
}
