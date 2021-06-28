package routers

import (
	"user/core/dependencies/services"

	"google.golang.org/grpc"
)

type RouteHandler interface {
	Register(server *grpc.Server, appDependencies services.App)
}
