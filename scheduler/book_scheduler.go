package scheduler

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookScheduler struct {
	db *mongo.Database
}

func NewBookScheduler(db *mongo.Database) *BookScheduler {
	return &BookScheduler{db: db}
}

func (s *BookScheduler) Start() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			s.checkOverdueBooks()
		}
	}()
}

func (s *BookScheduler) checkOverdueBooks() {
	ctx := context.Background()
	now := time.Now()

	filter := bson.M{
		"return_date": nil,
		"borrowed_date": bson.M{
			"$lt": now.Add(-14 * 24 * time.Hour),
		},
	}

	update := bson.M{
		"$set": bson.M{"status": "overdue"},
	}

	result, err := s.db.Collection("borrowed_books").UpdateMany(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating overdue books: %v", err)
		return
	}

	log.Printf("Updated %d overdue books", result.ModifiedCount)
}
