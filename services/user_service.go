package services

import (
	"context"

	"gc-buku/models"
	pb "gc-buku/proto"
	"gc-buku/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	db *mongo.Database
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := models.User{
		Username: req.User.Username,
		Password: req.User.Password,
	}

	result, err := s.db.Collection("users").InsertOne(ctx, user)
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

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	var user models.User
	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch user")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       user.ID.Hex(),
			Username: user.Username,
			Password: user.Password,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User
	err := s.db.Collection("users").FindOne(ctx, bson.M{
		"username": req.Username,
		"password": req.Password,
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "invalid credentials")
		}
		return nil, status.Errorf(codes.Internal, "failed to authenticate")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{
		Token: token,
		User: &pb.User{
			Id:       user.ID.Hex(),
			Username: user.Username,
		},
	}, nil
}
