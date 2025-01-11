package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"gc-buku/models"
	pb "gc-buku/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedBookServiceServer
	db *mongo.Database
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userCollection := s.db.Collection("users")
	user := models.User{
		Username: req.User.Username,
		Password: req.User.Password,
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	objectID := result.InsertedID.(primitive.ObjectID)

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:       objectID.Hex(),
			Username: user.Username,
			Password: user.Password,
		},
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userCollection := s.db.Collection("users")
	var user models.User

	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       user.ID.Hex(),
			Username: user.Username,
			Password: user.Password,
		},
	}, nil
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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookServiceServer(s, &server{db: db})

	log.Printf("gRPC Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
