package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/MaxPolarfox/tasks-rest/pkg/rest_client"
	tasksTypes "github.com/MaxPolarfox/tasks-rest/pkg/types"
	"github.com/MaxPolarfox/goTools/errors"
)


// RestTasksController handles GET, POST, DELETE, HEAD /rest/tasks calls
type RestTasksController struct {
	restTasksClient rest_client.Client
}

// NewRestTasksController creates a new RestTasksController
func NewRestTasksController(tasksClient rest_client.Client) *RestTasksController {
	return &RestTasksController{
		restTasksClient: tasksClient,
	}
}

// AddTask POST /tasks
func (c *RestTasksController) AddTask(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "RestTasksController.AddTask"

	body := tasksTypes.AddTaskReqBody{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		log.Println(metricName+".decode", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := c.restTasksClient.AddTask(ctx, body)
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
func (c *RestTasksController) GetAllTasks(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "RestTasksController.GetAllTasks"

	tasks, err := c.restTasksClient.GetAllTasks(ctx)
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
func (c *RestTasksController) DeleteTask(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	metricName := "RestTasksController.DeleteTask"

	params := httprouter.ParamsFromContext(ctx)
	taskID := params.ByName("id")

	if len(taskID) == 0 {
		errors.RespondWithError(rw, http.StatusBadRequest, "no task id")
	}

	err := c.restTasksClient.DeleteTask(ctx, taskID)
	if err != nil {
		log.Println(metricName+".tasksClient.DeleteTask", "err", err)
		errors.RespondWithError(rw, http.StatusNotFound, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}
