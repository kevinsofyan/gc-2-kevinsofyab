package main

import (
	"log"
	"os"

	"gc-buku/client/routes"

	pb "gc-buku/proto"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	grpcServer := os.Getenv("GRPC_SERVER")
	if grpcServer == "" {
		grpcServer = "server:50051"
	}

	// Setup gRPC connection
	conn, err := grpc.Dial(grpcServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewBookServiceClient(conn)

	routes.RegisterRoutes(e, client)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
