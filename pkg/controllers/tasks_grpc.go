package controllers

import (
	"encoding/json"
	"github.com/MaxPolarfox/gateway/pkg/types"
	"github.com/MaxPolarfox/goTools/errors"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"

	grpcTasksClient "github.com/MaxPolarfox/tasks/pkg/client"
)

// GrpcTasksController handles GET, POST, DELETE, HEAD /grpc/tasks calls
type GrpcTasksController struct {
	grpcTasksClient grpcTasksClient.Client
}

// NewGrpcTasksController creates a new GrpcTasksController
func NewGrpcTasksController(tasksClient grpcTasksClient.Client) *GrpcTasksController {
	return &GrpcTasksController{
		grpcTasksClient: tasksClient,
	}
}

// CreateTask POST /tasks
func (c *GrpcTasksController) CreateTask(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "GrpcTasksController.CreateTask"

	body := types.AddTaskReqBody{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		log.Println(metricName+".decode", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := c.grpcTasksClient.CreateTask(ctx, body.Data)
	if err != nil {
		log.Println(metricName, "err", err)
		errors.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	js, err := json.Marshal(res)
	if err != nil {
		log.Println(metricName+".Marshal", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(js)
}

// GetAllTasks GET /tasks
func (c *GrpcTasksController) GetTasks(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "GrpcTasksController.GetTasks"

	tasks, err := c.grpcTasksClient.GetTasks(ctx)
	if err != nil {
		log.Println(metricName, "err", err)
		errors.RespondWithError(rw, http.StatusBadRequest, err.Error())
		return
	}

	js, err := json.Marshal(tasks)
	if err != nil {
		log.Println(metricName+".Marshal", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
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
		errors.RespondWithError(rw, http.StatusBadRequest, "no task id")
	}

	err := c.grpcTasksClient.DeleteTask(ctx, taskID)
	if err != nil {
		log.Println(metricName+".tasksClient.DeleteTask", "err", err)
		errors.RespondWithError(rw, http.StatusNotFound, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}