package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes type is the container for the webRouter
type Routes struct {
	webRouter *mux.Router
}

// NewRoutes creates the Routes object that contains the webrouter to be referenced throughout the application
func NewRoutes() *Routes {
	return new(Routes)
}

// Begin is the main entry point for the web request router
func (r *Routes) Begin(port string) {
	r.webRouter = mux.NewRouter().StrictSlash(true)

	r.webRouter.HandleFunc("/login", SessionsHandlerLogin).Methods("POST")
	r.webRouter.HandleFunc("/logout", CheckLoginStatus(SessionsHandlerLogout)).Methods("POST")

	// No authentication required
	r.webRouter.HandleFunc("/users", UsersHandlerCreate).Methods("Post")
	r.webRouter.HandleFunc("/users", CheckLoginStatus(UsersHandlerGet)).Methods("Get")
	r.webRouter.HandleFunc("/users", CheckLoginStatus(UsersHandlerUpdate)).Methods("Put")
	r.webRouter.HandleFunc("/users", CheckLoginStatus(UsersHandlerDelete)).Methods("Delete")

	r.webRouter.HandleFunc("/tasks", CheckLoginStatus(TasksHandlerGetAll)).Methods("Get")
	r.webRouter.HandleFunc("/tasks", CheckLoginStatus(TasksHandlerCreate)).Methods("Post")
	r.webRouter.HandleFunc("/tasks/{taskid}", CheckLoginStatus(TasksHandlerGet)).Methods("Get")
	r.webRouter.HandleFunc("/tasks/{taskid}", CheckLoginStatus(TasksHandlerUpdate)).Methods("Put")
	r.webRouter.HandleFunc("/tasks/{taskid}", CheckLoginStatus(TasksHandlerDelete)).Methods("Delete")

	r.webRouter.HandleFunc("/tasks/{taskid}/reports", CheckLoginStatus(ReportsHandlerGetAll)).Methods("Get")
	r.webRouter.HandleFunc("/tasks/{taskid}/reports", CheckLoginStatus(ReportsHandlerCreate)).Methods("Post")
	r.webRouter.HandleFunc("/tasks/{taskid}/reports/{reportid}", CheckLoginStatus(ReportsHandlerGet)).Methods("Get")
	r.webRouter.HandleFunc("/tasks/{taskid}/reports/{reportid}", CheckLoginStatus(ReportsHandlerUpdate)).Methods("Put")
	r.webRouter.HandleFunc("/tasks/{taskid}/reports/{reportid}", CheckLoginStatus(ReportsHandlerDelete)).Methods("Delete")

	mlog.Info("Webserver up and running on port 3000.")
	mlog.Fatal(http.ListenAndServe(port, r.webRouter))
}
