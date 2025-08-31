package container

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lokker96/grpc_project/domain/service"
	"github.com/lokker96/grpc_project/infrastructure/persistence/postgres"
)

// Define the Container structure
// This is useful for setting up the internal container infrastructure
// and hide complexity from the main function
type Container struct {
	// ctx            context.Context // Context for the container
	ExplorerServer *service.ExploreServer
}

// NewContainer function creates and returns a new Container instance
func NewContainer() (*Container, error) {

	// This can be stored using docker secrets or 3rd party solution
	dbPassword := "testingPassword"

	// Setup the timezone for the database
	zone, _ := time.Now().Zone()

	// Prepare connection string for database connection
	dbDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		string(dbPassword),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		zone,
	)

	// Creatre new db connection using gorm
	dbConnection, err := NewDBConnection(dbDSN)
	if err != nil {
		return nil, fmt.Errorf("error on creating new db connection: %w", err)
	}

	// Build new explorer repository with the db connection created before
	explorerRepository := postgres.NewExplorerRepository(context.Background(), dbConnection)

	// Create the explorer server using the gRPC server code and attach the explorer
	// repository that implements the database routines for accessing the data using gorm
	explorerServer := service.NewExplorerServer(explorerRepository)

	// Create some dummy data - disable if you don't want it
	explorerServer.BuildDummyDataset()

	// Return a new Container instance with its explorer server
	return &Container{
		ExplorerServer: explorerServer,
	}, nil
}
