package gateway

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/MaxPolarfox/gateway/pkg/types"
	"github.com/MaxPolarfox/gateway/pkg/controllers"
	grpcTasksClient "github.com/MaxPolarfox/tasks/pkg/client"
	restTasksClient "github.com/MaxPolarfox/tasks-rest/pkg/rest_client"
)

type Service struct {
	Options types.Options
	Router *httprouter.Router
}

func NewService(options types.Options) *Service{

	// Clients
	grpcTasksClient := grpcTasksClient.NewTasksClient(options.Services.TasksGrpc)
	restTasksClient := restTasksClient.NewTasksClient(options.Services.TasksRest)

	// Controllers
	tasksGrpcController := controllers.NewGrpcTasksController(grpcTasksClient)
	tasksRestController := controllers.NewRestTasksController(restTasksClient)

	router := httprouter.New()

	// Routes
	router.HandlerFunc(http.MethodPost, "/grpc/tasks/", tasksGrpcController.CreateTask)
	router.HandlerFunc(http.MethodGet, "/grpc/tasks", tasksGrpcController.GetTasks)
	router.HandlerFunc(http.MethodDelete, "/grpc/tasks/:id", tasksGrpcController.DeleteTask)

	router.HandlerFunc(http.MethodPost, "/rest/tasks/", tasksRestController.CreateTask)
	router.HandlerFunc(http.MethodGet, "/rest/tasks", tasksRestController.CreateTask)
	router.HandlerFunc(http.MethodDelete, "/rest/tasks/:id", tasksRestController.DeleteTask)

	return &Service{
		Router: router,
		Options: options,
	}
}
