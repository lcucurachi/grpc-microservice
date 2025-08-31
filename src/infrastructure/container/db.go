package container

import (
	"github.com/lokker96/grpc_project/domain/entity"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Used to open a connection to a database with GORM and a posgres driver
func NewDBConnection(dsn string) (*gorm.DB, error) {

	// setup the entities used to build the tables in the DB
	EntityTypes := []interface{}{
		&entity.User{},
		&entity.Decision{},
	}

	// open the connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error on opening db connection: %w", err)
	}

	// tables are dropped just for demonstration purposes, dropping tables in production would cause data to be lost.
	for _, entityType := range EntityTypes {
		if err := db.Migrator().DropTable(entityType); err != nil {
			return nil, fmt.Errorf("error on dropping table: %w", err)
		}
	}

	// iterate through entities and use gorm tags to build tables
	for _, entityType := range EntityTypes {
		if err := db.AutoMigrate(entityType); err != nil {
			return nil, fmt.Errorf("error on auto migrating table: %w", err)
		}
	}

	return db, nil
}
