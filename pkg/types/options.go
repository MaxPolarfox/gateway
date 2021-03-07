package types

import goTools "github.com/MaxPolarfox/goTools/client"

type Options struct {
	Port int `json:"port"`
	ServiceName string `json:"serviceName"`
	Services Services         `json:"services"`
}

type Services struct {
	TasksRest  goTools.Options  `json:"tasks-rest"`
	TasksGrpc  goTools.Options  `json:"tasks-grpc"`
}





