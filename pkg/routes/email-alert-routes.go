package routes

import (
	"github.com/Ajesh8/UserTasks/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterEmailAlertRoutes = func(router *mux.Router) {
	router.HandleFunc("/emailalert", controllers.CreateEmail).Methods("POST")
}
