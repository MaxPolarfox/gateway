package controllers

import (
	"encoding/json"
	"github.com/MaxPolarfox/gateway/pkg/types"
	"log"
	"net/http"

	tasksClient "github.com/MaxPolarfox/tasks/pkg/client"
)

// CarsController handles GET, POST, DELETE, HEAD /toDoList calls
type TasksController struct {
	tasksClient tasksClient.Client
}

// NewTasksController creates a new TasksController
func NewTasksController(tasksClient tasksClient.Client) *TasksController {
	return &TasksController{
		tasksClient: tasksClient,
	}
}

// AddTask POST /toDoList
func (c *TasksController) AddTask(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "TasksController.AddTask"

	body := types.AddTaskReqBody{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		log.Println(metricName+".decode", "err", err)
		RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := c.tasksClient.AddTask(ctx, body.Data)
	if err != nil {
		log.Println(metricName, "err", err)
		RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	js, err := json.Marshal(res)
	if err != nil {
		log.Println(metricName+".Marshal", "err", err)
		RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(js)
}

type Error struct {
	Message string
}

func RespondWithError(rw http.ResponseWriter, statusCode int, message string) {

	response := Error{
		Message: message,
	}
	js, err := json.Marshal(&response)
	if err != nil {
		failedToMarshalError := Error{
			Message: err.Error(),
		}
		failedJS, _ := json.Marshal(&failedToMarshalError)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(failedJS)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(js)
}