package gateway

import (
	"github.com/MaxPolarfox/gateway/pkg/types"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/MaxPolarfox/gateway/pkg/controllers"
	tasksClient "github.com/MaxPolarfox/tasks/pkg/client"
)

type Service struct {
	Options types.Options
	Router *httprouter.Router
}

func NewService(options types.Options) *Service{

	// Clients
	tasksClient := tasksClient.NewTasksClient()

	// Controllers
	tasksController := controllers.NewTasksController(tasksClient)

	router := httprouter.New()

	// Routes
	router.HandlerFunc(http.MethodPost, "/tasks", tasksController.AddTask)
	router.HandlerFunc(http.MethodGet, "/tasks", tasksController.GetAllTasks)
	router.HandlerFunc(http.MethodDelete, "/tasks/:id", tasksController.DeleteTask)

	return &Service{
		Router: router,
		Options: options,
	}
}
