package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "gc-buku/proto"
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
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return s.userService.CreateUser(ctx, req)
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return s.userService.GetUser(ctx, req)
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.userService.Login(ctx, req)
}

func (s *server) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	return s.bookService.CreateBook(ctx, req)
}

func (s *server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	return s.bookService.GetBook(ctx, req)
}

func (s *server) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	return s.bookService.UpdateBook(ctx, req)
}

func (s *server) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	return s.bookService.DeleteBook(ctx, req)
}

func (s *server) BorrowBook(ctx context.Context, req *pb.BorrowBookRequest) (*pb.BorrowBookResponse, error) {
	return s.borrowService.BorrowBook(ctx, req)
}

func (s *server) ReturnBook(ctx context.Context, req *pb.ReturnBookRequest) (*pb.ReturnBookResponse, error) {
	return s.borrowService.ReturnBook(ctx, req)
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
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "book_management"
	}

	db, err := initDB(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	srv := &server{
		userService:   services.NewUserService(db),
		bookService:   services.NewBookService(db),
		borrowService: services.NewBorrowService(db),
	}

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
