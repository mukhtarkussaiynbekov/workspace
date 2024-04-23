package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mukhtarkv/workspace/api/todo/todoapp/v1beta1"
	"github.com/mukhtarkv/workspace/kit"
	"github.com/mukhtarkv/workspace/kit/config"
	"github.com/mukhtarkv/workspace/kit/id"
	"github.com/mukhtarkv/workspace/kit/log"
	db "github.com/mukhtarkv/workspace/kit/sql"
	"github.com/mukhtarkv/workspace/todo/todoapp"
	"github.com/mukhtarkv/workspace/todo/todoapp/assets"
	"github.com/mukhtarkv/workspace/todo/todoapp/postgres"
	"google.golang.org/grpc"
)

// Main function.
// Everything start from here.
func main() {
	podName := config.LookupEnv("POD_NAME", id.NewGenerator("todo-todoapp").Generate())
	ctx := context.Background()
	// Initiate a logger with pre-configuration for production and telemetry.
	l, err := log.New()
	if err != nil {
		// in case we cannot create the logger, the app should immediately stop.
		panic(err)
	}
	// Replace the global logger with the Service scoped log.
	log.ReplaceGlobal(l)

	// Initialize service
	// Mostly business logic initialization will be there
	storage, err := postgres.New()
	if err != nil {
		l.Fatal(ctx, err.Error())
	}

	// Apply migration scripts
	if err := db.Migrate(storage.DB, "todo", assets.SF); err != nil {
		l.Fatal(ctx, err.Error())
	}

	todoService := todoapp.New(storage)
	srv, err := newGrpcToDo(todoService)
	if err != nil {
		l.Fatal(ctx, err.Error())
	}

	// Initialise the foundation and start the service
	foundation, err := kit.NewFoundation("todoapp", kit.WithLogger(l))
	if err != nil {
		l.Fatal(ctx, err.Error())
	}

	// Register the GRPC Server.
	foundation.RegisterService(func(s *grpc.Server) {
		pb.RegisterToDoAppServer(s, srv)
	})

	// Register the Service Handler.
	foundation.RegisterServiceHandler(func(gw *runtime.ServeMux, conn *grpc.ClientConn) {
		if err := pb.RegisterToDoAppHandler(ctx, gw, conn); err != nil {
			l.Fatal(ctx, "fail registering gateway handler", log.Error(err))
		}
	})

	l.Info(ctx, "Starting service", log.String("pod.name", podName))

	// Start the service
	//
	// This service will be automatically configured to:
	// 1. Provide Observability information such as tracing, loging and metric
	// 2. Provide default /readyz and /healthz endpoint for rediness and liveness probe and profiling via /debug/pprof
	// 3. Setup for production setup
	// 4. Graceful shutdown
	if err := foundation.Serve(); err != nil {
		l.Error(ctx, "fail serving", log.Error(err))
	}
}
