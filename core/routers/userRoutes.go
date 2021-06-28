package routers

import (
	"user/core/dependencies/services"
	"user/core/handlers"

	"google.golang.org/grpc"
)

func NewUserRouter() RouteHandler {
	return userRoutes{}
}

type userRoutes struct {
}

func (router userRoutes) Register(server *grpc.Server, appDependencies services.App) {
	userHandler := handlers.NewUserHandler(appDependencies)
	userHandler.RegisterCrud(server)
}
