package main

import (
	"fmt"
	"log"
	"net"
	"user/core/dependencies/dependency"
	"user/core/dependencies/services"
	"user/core/routers"

	"google.golang.org/grpc"
)

type server struct {
	listener net.Listener
	grpc     *grpc.Server
	address  string
}

func (object *server) start() {
	fmt.Println("Server Listening...")
	if err := object.grpc.Serve(object.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func SetupApp(dependencies *services.App) (*server, error) {
	appDependencies := dependencies
	if appDependencies == nil {
		defaultDependencies, err := dependency.NewServiceLocator(nil)
		if err != nil {
			return nil, err
		}

		appDependencies = defaultDependencies
	}

	listener, err := net.Listen(
		appDependencies.Credentials.ServerTransportProtocol,
		appDependencies.Credentials.FullServerAddress(),
	)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return nil, err
	}

	grpcServer := grpc.NewServer()
	userRouter := routers.NewUserRouter()
	userRouter.Register(grpcServer, *appDependencies)

	return &server{
		listener: listener,
		grpc:     grpcServer,
		address:  appDependencies.Credentials.ServerAddress,
	}, nil
}

func main() {
	server, err := SetupApp(nil)
	if err != nil {
		log.Fatalf("Can't init app, %v", err)
		return
	}

	server.start()
}
