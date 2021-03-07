package gateway

import (
	"github.com/MaxPolarfox/gateway/pkg/types"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/MaxPolarfox/gateway/pkg/controllers"
	grpcTasksClient "github.com/MaxPolarfox/tasks/pkg/grpc_client"
)

type Service struct {
	Options types.Options
	Router *httprouter.Router
}

func NewService(options types.Options) *Service{

	// Clients
	grpcTasksClient := grpcTasksClient.NewTasksClient()

	// Controllers
	tasksController := controllers.NewGrpcTasksController(grpcTasksClient)

	router := httprouter.New()

	// Routes
	router.HandlerFunc(http.MethodPost, "/grpc/tasks/", tasksController.AddTask)
	router.HandlerFunc(http.MethodGet, "/grpc/tasks", tasksController.GetAllTasks)
	router.HandlerFunc(http.MethodDelete, "/grpc/tasks/:id", tasksController.DeleteTask)

	return &Service{
		Router: router,
		Options: options,
	}
}
