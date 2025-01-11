package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	Author        string             `bson:"author"`
	PublishedDate time.Time          `bson:"published_date"`
	Status        string             `bson:"status"`
	UserID        primitive.ObjectID `bson:"user_id,omitempty"`
}
