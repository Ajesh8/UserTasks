package routes

import (
	"github.com/Ajesh8/UserTasks/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterLoginRoutes = func(router *mux.Router) {
	router.HandleFunc("/signup", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.Signin).Methods("POST")
}
