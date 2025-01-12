package services

import (
	"context"
	"time"

	"gc-buku/models"
	pb "gc-buku/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookService struct {
	db *mongo.Database
}

func NewBookService(db *mongo.Database) *BookService {
	return &BookService{db: db}
}

func (s *BookService) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	book := models.Book{
		Title:         req.Book.Title,
		Author:        req.Book.Author,
		PublishedDate: time.Now(),
		Status:        "available",
	}

	result, err := s.db.Collection("books").InsertOne(ctx, book)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create book: %v", err)
	}

	objectID := result.InsertedID.(primitive.ObjectID)
	return &pb.CreateBookResponse{
		Book: &pb.Book{
			Id:            objectID.Hex(),
			Title:         book.Title,
			Author:        book.Author,
			PublishedDate: book.PublishedDate.Format(time.RFC3339),
			Status:        book.Status,
		},
	}, nil
}

func (s *BookService) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	var book models.Book
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID")
	}

	err = s.db.Collection("books").FindOne(ctx, bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "book not found")
	}

	return &pb.GetBookResponse{
		Book: &pb.Book{
			Id:            book.ID.Hex(),
			Title:         book.Title,
			Author:        book.Author,
			PublishedDate: book.PublishedDate.Format(time.RFC3339),
			Status:        book.Status,
			UserId:        book.UserID.Hex(),
		},
	}, nil
}

func (s *BookService) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Book.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID")
	}

	update := bson.M{
		"$set": bson.M{
			"title":          req.Book.Title,
			"author":         req.Book.Author,
			"published_date": req.Book.PublishedDate,
			"status":         req.Book.Status,
		},
	}

	_, err = s.db.Collection("books").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update book")
	}

	return &pb.UpdateBookResponse{Book: req.Book}, nil
}

func (s *BookService) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID")
	}

	_, err = s.db.Collection("books").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete book")
	}

	return &pb.DeleteBookResponse{Id: req.Id}, nil
}
