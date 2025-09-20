package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoDB_URI       string
	MongoDB_URI_DB1   string
	MongoDB_URI_DB2   string
	DatabaseName      string
	DatabaseName_DB1 string
	DatabaseName_DB2 string
	RedisAddr         string
	ElasticsearchURL  string
	ServerPort        string
	JWTSecret         string
}

func Load() *Config {
	return &Config{
		MongoDB_URI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB_URI_DB1:   getEnv("MONGO_URI_DB1", "mongodb://localhost:27017"),
		MongoDB_URI_DB2:   getEnv("MONGO_URI_DB2", "mongodb://localhost:27017"),
		DatabaseName:      getEnv("MONGO_DATABASE", "bookmyshow"),
		DatabaseName_DB1:  getEnv("MONGO_DATABASE_DB1", "bookmyshow_db1"),
		DatabaseName_DB2:  getEnv("MONGO_DATABASE_DB2", "bookmyshow_db2"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		ElasticsearchURL:  getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ConnectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB")
	return client, nil
}

func ConnectMongoDB1(cfg *Config) (*mongo.Client, error) {
	return ConnectMongo(cfg.MongoDB_URI_DB1)
}

func ConnectMongoDB2(cfg *Config) (*mongo.Client, error) {
	return ConnectMongo(cfg.MongoDB_URI_DB2)
}

func ConnectRedis(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}

	return client
}

func ConnectElasticsearch(url string) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Warning: Failed to create Elasticsearch client: %v", err)
		return nil
	}

	res, err := client.Info()
	if err != nil {
		log.Printf("Warning: Failed to connect to Elasticsearch: %v", err)
		return nil
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Warning: Elasticsearch returned error: %s", res.String())
		return nil
	}

	log.Println("Connected to Elasticsearch")
	return client
}
