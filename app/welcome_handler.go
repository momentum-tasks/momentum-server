package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
)

// UsersHandler is a debug webpage that returns the user with the id of 1
// This should really be removed once main functionality is complete
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := GetUserByID(1)
	j, _ := json.MarshalIndent(u, "", "    ")
	fmt.Fprint(w, string(j))
}

// CreateUsersHandler is a debug webpage that creates a sample user with a task and report
// This should be removed once additional testing methods are created
func CreateUsersHandler(w http.ResponseWriter, r *http.Request) {
	CreateUser("bdylan", "bob@example.com", "password")
	u, _ := GetUserByID(1)
	CreateTask(u, "Test Task", "This is only a test", pq.NullTime{}, 1, false)
	tasks := GetTasksByUser(u)
	CreateReport(&tasks[0], 1, "This is my first report")
}

// SessionsHandler is a debug webpage that returns the sessions of the user with the id of 1
// This should really be removed once main functionality is complete
func SessionsHandler(w http.ResponseWriter, r *http.Request) {
	sessions := GetSessionsByUser(1)
	j, _ := json.MarshalIndent(sessions, "", "    ")
	fmt.Fprint(w, string(j))
}
