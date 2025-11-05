package repository

import (
	"context"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ConfigRepository struct {
	collection *mongo.Collection
}

func NewConfigRepository(client *mongo.Client, dbName string) *ConfigRepository {
	collection := client.Database(dbName).Collection("config")
	return &ConfigRepository{collection: collection}
}

func (r *ConfigRepository) Create(config *models.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, config)
	return err
}

func (r *ConfigRepository) FindAll() ([]models.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []models.Config
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ConfigRepository) FindByDBType(dbType string) ([]models.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"dbType": dbType})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []models.Config
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ConfigRepository) FindByDBTypeAndHashID(dbType, hashID string) (*models.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	combinedID := dbType + "_" + hashID
	var config models.Config
	err := r.collection.FindOne(ctx, bson.M{"_id": combinedID}).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigRepository) Update(config *models.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": config.ID}, bson.M{"$set": config})
	return err
}

func (r *ConfigRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
