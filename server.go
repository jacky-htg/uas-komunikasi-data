package main

import (
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"uas-komdat/internal/config"
	"uas-komdat/internal/middleware"
	"uas-komdat/internal/pkg/db/postgres"
	"uas-komdat/internal/pkg/storage"
	"uas-komdat/internal/route"
)

const defaultPort = "8000"

func main() {
	// lookup and setup env
	if _, ok := os.LookupEnv("PORT"); !ok {
		config.Setup(".env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// init log
	log := log.New(os.Stdout, "LMS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// create postgres database connection
	db, err := postgres.Open()
	if err != nil {
		log.Fatalf("connecting to db: %v", err)
		return
	}
	log.Print("connecting to postgresql database")

	defer db.Close()

	minioClient, err := storage.Connection()
	if err != nil {
		log.Fatalf("connecting to minio: %v", err)
		return
	}

	// listen tcp port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	authInterceptor := middleware.Context{}
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			authInterceptor.Unary(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			authInterceptor.Stream(),
		)),
	}

	grpcServer := grpc.NewServer(serverOptions...)

	// routing grpc services
	route.GrpcRoute(grpcServer, db, log, minioClient)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
		return
	}
	log.Print("serve grpc on port: " + port)
}
