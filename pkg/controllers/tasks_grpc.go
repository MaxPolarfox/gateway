package controllers

import (
	"encoding/json"
	"github.com/MaxPolarfox/gateway/pkg/types"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"

	grpcTasksClient "github.com/MaxPolarfox/tasks/pkg/grpc_client"
)

// CarsController handles GET, POST, DELETE, HEAD /toDoList calls
type GrpcTasksController struct {
	grpcTasksClient grpcTasksClient.Client
}

// NewTasksController creates a new TasksController
func NewGrpcTasksController(tasksClient grpcTasksClient.Client) *GrpcTasksController {
	return &GrpcTasksController{
		grpcTasksClient: tasksClient,
	}
}

// AddTask POST /tasks
func (c *GrpcTasksController) AddTask(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "GrpcTasksController.AddTask"

	body := types.AddTaskReqBody{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		log.Println(metricName+".decode", "err", err)
		RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := c.grpcTasksClient.AddTask(ctx, body.Data)
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

// GetAllTasks GET /tasks
func (c *GrpcTasksController) GetAllTasks(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "GrpcTasksController.GetAllTasks"

	tasks, err := c.grpcTasksClient.GetAllTasks(ctx)
	if err != nil {
		log.Println(metricName, "err", err)
		RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	js, err := json.Marshal(tasks)
	if err != nil {
		log.Println(metricName+".Marshal", "err", err)
		RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(js)
}

// DeleteTask DELETE /tasks/:id
func (c *GrpcTasksController) DeleteTask(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "GrpcTasksController.DeleteTask"

	params := httprouter.ParamsFromContext(ctx)

	taskID := params.ByName("id")

	if len(taskID) == 0 {
		RespondWithError(rw, http.StatusBadRequest, "no task id")
	}

	err := c.grpcTasksClient.DeleteTask(ctx, taskID)
	if err != nil {
		log.Println(metricName+".tasksClient.DeleteTask", "err", err)
		RespondWithError(rw, http.StatusNotFound, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
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