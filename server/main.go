package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "gc-buku/proto"
	"gc-buku/scheduler"
	"gc-buku/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBookServiceServer
	userService   *services.UserService
	bookService   *services.BookService
	borrowService *services.BorrowService
	bookScheduler *scheduler.BookScheduler
}

func initDB(mongoURI, dbName string) (*mongo.Database, error) {
	ctx := context.TODO()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping: %v", err)
	}

	db := client.Database(dbName)

	collections := []string{"users", "books", "borrowed_books"}
	for _, col := range collections {
		if err := db.CreateCollection(ctx, col); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return nil, fmt.Errorf("failed to create collection %s: %v", col, err)
			}
		}
	}

	log.Printf("Connected to MongoDB: %s", dbName)
	return db, nil
}

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	db, err := initDB(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize services
	srv := &server{
		userService:   services.NewUserService(db),
		bookService:   services.NewBookService(db),
		borrowService: services.NewBorrowService(db),
		bookScheduler: scheduler.NewBookScheduler(db),
	}

	// Start scheduler
	srv.bookScheduler.Start()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookServiceServer(s, srv)

	log.Printf("gRPC Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
