package repository

import (
	"context"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TheatreRepository struct {
	collection *mongo.Collection
}

func NewTheatreRepository(client *mongo.Client, dbName string) *TheatreRepository {
	collection := client.Database(dbName).Collection("theatres")
	return &TheatreRepository{collection: collection}
}

func (r *TheatreRepository) Create(theatre *models.Theatre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	theatre.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, theatre)
	return err
}

func (r *TheatreRepository) FindAll() ([]models.Theatre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var theatres []models.Theatre
	if err := cursor.All(ctx, &theatres); err != nil {
		return nil, err
	}
	return theatres, nil
}

func (r *TheatreRepository) FindByID(id string) (*models.Theatre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var theatre models.Theatre
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&theatre)
	if err != nil {
		return nil, err
	}
	return &theatre, nil
}

func (r *TheatreRepository) Update(theatre *models.Theatre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": theatre.ID}, bson.M{"$set": theatre})
	return err
}

func (r *TheatreRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}