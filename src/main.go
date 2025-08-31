package main

import (
	"log"
	"net"

	"github.com/lokker96/grpc_project/infrastructure/container"
	ep "github.com/lokker96/grpc_project/infrastructure/proto/explore"
	"google.golang.org/grpc"
)

func main() {
	// Building the application's container
	c, err := container.NewContainer() // Create a new container instance
	if err != nil {
		log.Fatal(err) // If there's an error, log it and terminate
	}

	// Setup a listener on port 9001
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal("failed to listen: ", err.Error())
	}

	// Create new gRPC server and set the service responsable for responding
	grpcServer := grpc.NewServer()
	ep.RegisterExploreServiceServer(grpcServer, c.ExplorerServer)

	// Listen to new gRPC calls
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}
