package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BorrowedBook struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	BookID       primitive.ObjectID `bson:"book_id"`
	UserID       primitive.ObjectID `bson:"user_id"`
	BorrowedDate time.Time          `bson:"borrowed_date"`
	ReturnDate   *time.Time         `bson:"return_date,omitempty"`
}
