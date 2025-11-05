package repository

import (
	"context"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MovieRepository struct {
	collection *mongo.Collection
}

func NewMovieRepository(client *mongo.Client, dbName string) *MovieRepository {
	collection := client.Database(dbName).Collection("movies")
	return &MovieRepository{collection: collection}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movie.ID = primitive.NewObjectID()
	movie.LastUpdate = time.Now()
	_, err := r.collection.InsertOne(ctx, movie)
	return err
}

func (r *MovieRepository) FindAll() ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []models.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepository) FindByID(id string) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var movie models.Movie
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) FindAfterLastUpdatedTime(lastUpdatedTime string) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t, err := time.Parse(time.RFC3339, lastUpdatedTime)
	if err != nil {
		t = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"lastUpdated": bson.M{"$gt": t}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []models.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepository) Update(movie *models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": movie.ID}, bson.M{"$set": movie})
	return err
}

func (r *MovieRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}