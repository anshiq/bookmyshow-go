package repository

import (
	"context"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShowRepository struct {
	collection *mongo.Collection
}

func NewShowRepository(client *mongo.Client, dbName string) *ShowRepository {
	collection := client.Database(dbName).Collection("shows")
	return &ShowRepository{collection: collection}
}

func (r *ShowRepository) Create(show *models.Show) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	show.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, show)
	return err
}

func (r *ShowRepository) FindAll() ([]models.Show, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shows []models.Show
	if err := cursor.All(ctx, &shows); err != nil {
		return nil, err
	}
	return shows, nil
}

func (r *ShowRepository) FindByID(id string) (*models.Show, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var show models.Show
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&show)
	if err != nil {
		return nil, err
	}
	return &show, nil
}

func (r *ShowRepository) FindByMovieID(movieID string) ([]models.Show, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"movieId": movieID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shows []models.Show
	if err := cursor.All(ctx, &shows); err != nil {
		return nil, err
	}
	return shows, nil
}

func (r *ShowRepository) Update(show *models.Show) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": show.ID}, bson.M{"$set": show})
	return err
}

func (r *ShowRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}