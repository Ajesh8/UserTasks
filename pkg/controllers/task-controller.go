package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ajesh8/UserTasks/pkg/models"
	"github.com/Ajesh8/UserTasks/pkg/utils"
	"github.com/gorilla/mux"
)

type Message struct {
	Msg string `json:"Message"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	authorized, user := utils.CheckUserAuthorisation(w, r)
	if !authorized {
		return
	}

	CreateTask := &models.Task{}
	err := utils.ParseBody(r, CreateTask)
	if err != nil {
		log.Printf("Could not parse body:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := CreateTask.CreateTask(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(&Message{
			Msg: err.Error(),
		})
		w.Write(res)
		return
	}
	res, err := json.Marshal(task)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	authorized, _ := utils.CheckUserAuthorisation(w, r)
	if !authorized {
		return
	}
	newTasks := models.GetAllTask()
	res, err := json.Marshal(newTasks)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func SearchTask(w http.ResponseWriter, r *http.Request) {
	authorized, _ := utils.CheckUserAuthorisation(w, r)
	if !authorized {
		return
	}
	searchText, ok := r.URL.Query()["title"]
	if !ok || len(searchText[0]) < 1 {
		log.Println("Url Param 'title' is missing")
		return
	}
	searchTasks := models.SearchTask(searchText[0])
	res, err := json.Marshal(searchTasks)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	authorized, _ := utils.CheckUserAuthorisation(w, r)
	if !authorized {
		return
	}
	vars := mux.Vars(r)
	taskId := vars["taskId"]
	ID, err := strconv.ParseInt(taskId, 0, 0)
	if err != nil {
		log.Println("Error while parsing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	status, ok := r.URL.Query()["status"]
	if !ok || len(status[0]) < 1 {
		log.Println("Url Param 'status' is missing")
		return
	}
	task := models.UpdateTaskStatus(ID, status[0])
	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res, err := json.Marshal(task)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func FilterTask(w http.ResponseWriter, r *http.Request) {
	authorized, _ := utils.CheckUserAuthorisation(w, r)
	if !authorized {
		return
	}
	dueText, ok := r.URL.Query()["due"]
	if !ok || len(dueText[0]) < 1 {
		log.Println("Url Param 'due' is missing")
		return
	}
	filteredTasks := models.FilterTask(dueText[0])
	res, err := json.Marshal(filteredTasks)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateEmail(w http.ResponseWriter, r *http.Request) {
	authorized, _ := utils.CheckUserAuthorisation(w, r)
	if !authorized {
		return
	}
	createEmailAlert := &models.EmailAlert{}
	err := utils.ParseBody(r, createEmailAlert)
	if err != nil {
		log.Printf("Could not parse body:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	emailAlert, err := createEmailAlert.CreateEmailAlert()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(&Message{
			Msg: err.Error(),
		})
		w.Write(res)
		return
	}
	res, err := json.Marshal(emailAlert)
	if err != nil {
		log.Printf("Marshalling error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
