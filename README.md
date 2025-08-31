# gRPC microservice in Golang

## Contents
1. [Introduction](#Introduction)
2. [Implementation Overview](#implementation-overview)
3. [Assumptions](#assumptions)
4. [Codebase Structure](#codebase-structure) 
5. [Run the code](#run-the-code)
6. [Generate gRPC Code](#generate-grpc-code)
7. [Generate the mocks](#generate-the-mocks)


## Introduction
The main goal of this project is to implement a gRPC microservice that responds to the 
protocol buffer definitions found at 'src/infrastructure/proto/explore/explore-service.proto'.

These definitions are used by an explorer service which is used to return data for decisions made by the users
in the context of a dating app. When a user likes another user a 'decision' is made. These decisions can
be overwritten too.

## Implementation Overview
The project has been implemented following this structure:

- Two containers: postgres-db, explorer-server.

- Database: PostgreSQL
  I've also used an ORM library (gorm) to simplify my operations with the database.
  Please look at 'src/domain/entities' files, each struct represents a database table.

- Codebase: I'm following a Domain Driven Design approach (domain, services, repositories)
  to structure the code.

- Errors: I kept a simple approach using the standard 'log' and 'fmt' libraries although I personally
  prefer using an external library that allows you to follow the errors more easily by providing a
  stacktrace which can be stored in the server logs.
  (I'm referencing this stacktrace library if someone is interested https://github.com/palantir/stacktrace)

- Dummy Data: I've added a routine to build some dummy data to play with the gRPC methods more easily
  but feel free to comment it out. The routine 'BuildDummyDataset()' is called inside: src/infrastructure/container.go

- gRPC Endpoints: I've implemented the routines inside 'src/infrastructure/domain/service/explorer_server.go'

- gRPC Client: I've implemented a client that can be used to call the routines on the server. Run this locally.
  I've used this one to test the application.

- Testing: I have written some unit tests using a mocking library called mockery (https://github.com/vektra/mockery)
  These unit tests would only test the domain/application layer. In order to test the system as a whole and actually check that the DB implementation works I would have to implement functional tests too.

- Pagination: I have avoided implementing it as this is a demo project. Pagination is a tricky topic because it depends on how
  does the application display the results. I would definitely try and use a timestamp approach if possible but being mindful about how the decisions are stored in the database.


## Assumptions
My main assumption in this project is that when a decision is made by a user (a like), only 1 row is created to represent this decision in the database. If the user decides to change their mind then we update this row. This way we always have 1 row per 'author_id' and 'recipient_id' pair, and another one for the opposite situation. 

Users that have decided on thousands of users might be problematic, especially in the way I've written the code because too much data is pulled and processed to find new likes so having dedicated SQL query will be a better method than processing everything by code.

The problem about the amount of rows in the decisions table can be improved by creating a two key index on the 'author_id' and 'recipient_id' columns since we can have only 1 record per pair, this way we can speed up the SQL query time. Of course there are other methods too (like region segregation) but this is a topic for a different discussion.


## Codebase Structure
Here I explain what the folders mean.

-- src - All the source code for the exercise

-- src/client - the client example code that can be used to make gRPC calls to the server

-- src/domain - this folder would contain the business logic but because we want to keep things simple I've
   only implemented entities (datatabase tables mapped to a structure), services (explorer service server) and custom erors. Ideally here we would have the business logic that isn't aware of the underlying implementation like a postgreSQL database or
   caching tools and so on. Repositories would describe interfaces that data storage solutions implement while services would 
   call these interfaces to provide some functionality for the endpoints. Then we would inject dependencies for database and other implemetations in the services inside the container by building services with specific repositories and so on. This would be too much for this example so I kept things simple but I wanted to give an idea of what a real microservice might look like.

-- src/domain/entities - the entities used in the database and in the business logic

-- src/domain/service - the services that we use to provide business logic

-- src/infrastructure - contains code that setups the microservice and implements the infrastructure, like the database

-- src/infrastructure/container - code that builds a container by initialising the explorer server, db and repositories

-- src/infrastructure/proto - contains the gRPC proto definitions

-- src/infrastructure/persistence/postgres - implements (DDD repository) methods to query the 
   postgreSQL database using gorm


## Run the code
First check if you want to disable the dummy data function call, then use docker compose to build and run the stack.

    docker compose build
    docker compose up

Once the explorer-server container is up, move into the client folder in 'src/client' and run it using 'go run client.go' to test the client code. Check the code to see what test cases I'm running.


## Regenerate gRPC code
If you need to update the proto code, then use the first command below to update the path env variable and then run the last command to regenerare the gRPC code.

    export PATH="$PATH:$(go env GOPATH)/bin"

    protoc \
        -I=$PWD/src/infrastructure/proto/explore \
        --go_out=$PWD/src/infrastructure/proto \
        --go-grpc_out=$PWD/src/infrastructure/proto \
        --go-grpc_opt=Mexplore-service.proto=./explore \
        --go_opt=Mexplore-service.proto=./explore \
        $PWD/src/infrastructure/proto/explore/explore-service.proto


## Generate the mocks
First install the mockery library (https://vektra.github.io/mockery/latest/installation/) and then execute the binary (https://vektra.github.io/mockery/latest/running/).
The config file in the src folder will take care of configuring which interfaces need to be mocked.

The mocks will be generated in 'src/mocks'.
