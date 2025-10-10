package repository

import (
	"context"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScreenRepository struct {
	collection *mongo.Collection
}

func NewScreenRepository(client *mongo.Client, dbName string) *ScreenRepository {
	collection := client.Database(dbName).Collection("screens")
	return &ScreenRepository{collection: collection}
}

func (r *ScreenRepository) Create(screen *models.Screen) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	screen.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, screen)
	return err
}

func (r *ScreenRepository) FindAll() ([]models.Screen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var screens []models.Screen
	if err := cursor.All(ctx, &screens); err != nil {
		return nil, err
	}
	return screens, nil
}

func (r *ScreenRepository) FindByID(id string) (*models.Screen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var screen models.Screen
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&screen)
	if err != nil {
		return nil, err
	}
	return &screen, nil
}

func (r *ScreenRepository) FindByTheatreID(theatreID string) ([]models.Screen, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"theatreId": theatreID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var screens []models.Screen
	if err := cursor.All(ctx, &screens); err != nil {
		return nil, err
	}
	return screens, nil
}

func (r *ScreenRepository) Update(screen *models.Screen) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": screen.ID}, bson.M{"$set": screen})
	return err
}

func (r *ScreenRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}