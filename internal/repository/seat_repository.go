package repository

import (
	"context"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatRepository struct {
	collection *mongo.Collection
}

func NewSeatRepository(client *mongo.Client, dbName string) *SeatRepository {
	collection := client.Database(dbName).Collection("seats")
	return &SeatRepository{collection: collection}
}

func (r *SeatRepository) Create(seat *models.Seat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	seat.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, seat)
	return err
}

func (r *SeatRepository) FindByNumberAndShowID(showID string, seatNumber int) (*models.Seat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"showId": showID, "seatNumber": seatNumber})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var seat models.Seat
	if cursor.Next(ctx) {
		err = cursor.Decode(&seat)
		if err != nil {
			return nil, err
		}
		return &seat, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (r *SeatRepository) Update(seat *models.Seat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": seat.ID}, bson.M{"$set": seat})
	return err
}

func (r *SeatRepository) UpdateVersion(showID string, seatNumber int, currentVersion int) (*models.Seat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"showId": showID, "seatNumber": seatNumber, "version": currentVersion}
	update := bson.M{"$inc": bson.M{"version": 2}, "$set": bson.M{"bookingStatus": true}}

	var seat models.Seat
	err := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&seat)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepository) FindAllByShowID(showID string) ([]models.Seat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"showId": showID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var seats []models.Seat
	if err := cursor.All(ctx, &seats); err != nil {
		return nil, err
	}
	return seats, nil
}
