package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// ReportsHandlerGetAll returns all reports associated with a given task
func ReportsHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	taskid, err := strconv.Atoi(mux.Vars(r)["taskid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		for _, t := range user.Tasks {
			if t.ID == taskid {
				j, _ := json.Marshal(t.Reports)
				fmt.Fprint(w, string(j))
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}

func ReportsHandlerCreate(w http.ResponseWriter, r *http.Request) {
}

// ReportsHandlerGet returns a single report associated with a given task
func ReportsHandlerGet(w http.ResponseWriter, r *http.Request) {
	taskid, err := strconv.Atoi(mux.Vars(r)["taskid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	reportid, err := strconv.Atoi(mux.Vars(r)["reportid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		for _, t := range user.Tasks {
			if t.ID == taskid {
				for _, report := range t.Reports {
					if report.ID == reportid {
						j, _ := json.Marshal(report)
						fmt.Fprint(w, string(j))
						return
					}
				}
				http.Error(w, "400 Bad Request", http.StatusBadRequest)
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}

func ReportsHandlerUpdate(w http.ResponseWriter, r *http.Request) {
}

func ReportsHandlerDelete(w http.ResponseWriter, r *http.Request) {
}
