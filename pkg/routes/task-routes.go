package routes

import (
	"github.com/Ajesh8/UserTasks/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterTasksRoutes = func(router *mux.Router) {
	router.HandleFunc("/task", controllers.CreateTask).Methods("POST")
	router.HandleFunc("/task", controllers.GetTask).Methods("GET")
	router.HandleFunc("/task/search", controllers.SearchTask).Methods("GET")
	router.HandleFunc("/task/{taskId}/update", controllers.UpdateTaskStatus).Methods("POST")
	router.HandleFunc("/task/filter", controllers.FilterTask).Methods("GET")

}
