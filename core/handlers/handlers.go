package handlers

import "google.golang.org/grpc"

type AppHandler interface {
	RegisterCrud(server *grpc.Server)
}
