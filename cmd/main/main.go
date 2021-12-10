package main

import (
	"log"
	"net/http"

	"github.com/Ajesh8/UserTasks/pkg/models"
	"github.com/Ajesh8/UserTasks/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterTasksRoutes(r)
	routes.RegisterLoginRoutes(r)
	routes.RegisterEmailAlertRoutes(r)
	http.Handle("/", r)
	log.Print("starting microservice")
	go models.AlertJob()
	log.Fatal(http.ListenAndServe("localhost:50000", r))
}
