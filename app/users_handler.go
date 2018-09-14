package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

// UsersHandlerGet returns all the information about the curent user
// This includes all their tasks and reports
func UsersHandlerGet(w http.ResponseWriter, r *http.Request) {
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		j, _ := json.Marshal(user)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(j))
	}
}

// UsersHandlerCreate is a POST endpoint to create a new user
// Requires a JSON body containing the username, password, and email
func UsersHandlerCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u struct {
		Username string
		Password string
		Email    string
	}
	err := decoder.Decode(&u)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if u.Username != "" && u.Password != "" && u.Email != "" {
		err = CreateUser(u.Username, u.Email, u.Password)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}
	}
}

func UsersHandlerUpdate(w http.ResponseWriter, r *http.Request) {
}

func UsersHandlerDelete(w http.ResponseWriter, r *http.Request) {
}
