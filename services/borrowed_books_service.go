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

type BorrowService struct {
	db *mongo.Database
}

func NewBorrowService(db *mongo.Database) *BorrowService {
	return &BorrowService{db: db}
}

func (s *BorrowService) BorrowBook(ctx context.Context, req *pb.BorrowBookRequest) (*pb.BorrowBookResponse, error) {
	bookID, err := primitive.ObjectIDFromHex(req.BorrowedBook.BookId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID")
	}

	userID, err := primitive.ObjectIDFromHex(req.BorrowedBook.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID")
	}

	// Start session for transaction
	session, err := s.db.Client().StartSession()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction")
	}
	defer session.EndSession(ctx)

	// Transaction
	result, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Check if book is available
		var book models.Book
		err := s.db.Collection("books").FindOne(ctx, bson.M{
			"_id":    bookID,
			"status": "available",
		}).Decode(&book)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, status.Errorf(codes.NotFound, "book not available")
			}
			return nil, status.Errorf(codes.Internal, "failed to fetch book")
		}

		// Create borrow record
		borrowedBook := models.BorrowedBook{
			BookID:       bookID,
			UserID:       userID,
			BorrowedDate: time.Now(),
		}

		borrowResult, err := s.db.Collection("borrowed_books").InsertOne(ctx, borrowedBook)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create borrow record")
		}

		// Update book status
		_, err = s.db.Collection("books").UpdateOne(
			ctx,
			bson.M{"_id": bookID},
			bson.M{"$set": bson.M{
				"status":  "borrowed",
				"user_id": userID,
			}},
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update book status")
		}

		return borrowResult.InsertedID, nil
	})

	if err != nil {
		return nil, err
	}

	objectID := result.(primitive.ObjectID)
	return &pb.BorrowBookResponse{
		BorrowedBook: &pb.BorrowedBook{
			Id:           objectID.Hex(),
			BookId:       bookID.Hex(),
			UserId:       userID.Hex(),
			BorrowedDate: time.Now().Format(time.RFC3339),
		},
	}, nil
}

func (s *BorrowService) ReturnBook(ctx context.Context, req *pb.ReturnBookRequest) (*pb.ReturnBookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid borrow ID")
	}

	session, err := s.db.Client().StartSession()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction")
	}
	defer session.EndSession(ctx)

	var borrowedBook models.BorrowedBook
	returnTime := time.Now()

	result, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Find and update borrowed book
		err := s.db.Collection("borrowed_books").FindOneAndUpdate(
			ctx,
			bson.M{
				"_id":         objectID,
				"return_date": nil,
			},
			bson.M{"$set": bson.M{"return_date": returnTime}},
		).Decode(&borrowedBook)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, status.Errorf(codes.NotFound, "borrow record not found or already returned")
			}
			return nil, status.Errorf(codes.Internal, "failed to update borrow record")
		}

		// Update book status
		_, err = s.db.Collection("books").UpdateOne(
			ctx,
			bson.M{"_id": borrowedBook.BookID},
			bson.M{
				"$set": bson.M{
					"status":  "available",
					"user_id": nil,
				},
			},
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update book status")
		}

		return &borrowedBook, nil
	})

	if err != nil {
		return nil, err
	}

	borrowedBook = *result.(*models.BorrowedBook)
	return &pb.ReturnBookResponse{
		BorrowedBook: &pb.BorrowedBook{
			Id:           borrowedBook.ID.Hex(),
			BookId:       borrowedBook.BookID.Hex(),
			UserId:       borrowedBook.UserID.Hex(),
			BorrowedDate: borrowedBook.BorrowedDate.Format(time.RFC3339),
			ReturnDate:   returnTime.Format(time.RFC3339),
		},
	}, nil
}
