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
			if t.Priority == taskid {
				j, _ := json.Marshal(t.Reports)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(j))
				return
			}
		}
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}

// ReportsHandlerCreate creates a report for the specified task, followed by a defrag of the task's reports
// Requires a JSON body with the report's sequence, and id
func ReportsHandlerCreate(w http.ResponseWriter, r *http.Request) {
	taskid, err := strconv.Atoi(mux.Vars(r)["taskid"])
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if rv := context.Get(r, UserContext); rv != nil {
		user := rv.(*User)
		decoder := json.NewDecoder(r.Body)
		var newReport struct {
			Sequence    int
			Description string
		}
		err := decoder.Decode(&newReport)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, t := range user.Tasks {
			if t.Priority == taskid {
				if newReport.Description != "" && newReport.Sequence > 0 {
					err = CreateReport(&t, newReport.Sequence, newReport.Description)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					err = DefragReports(&t)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					fmt.Fprint(w, "Successfully created report")
				}
				return
			}
		}
	}
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
			if t.Priority == taskid {
				for _, report := range t.Reports {
					if report.Sequence == reportid {
						j, _ := json.Marshal(report)
						w.Header().Set("Content-Type", "application/json")
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

// ReportsHandlerDelete deletes the report from the specified task, and then runs a defrag on the reports list for that task
func ReportsHandlerDelete(w http.ResponseWriter, r *http.Request) {
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
			if t.Priority == taskid {
				for _, report := range t.Reports {
					if report.Sequence == reportid {
						err = report.Delete()
						if err != nil {
							http.Error(w, "400 Bad Request", http.StatusBadRequest)
							return
						}
						err = DefragReports(&t)
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
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
}
