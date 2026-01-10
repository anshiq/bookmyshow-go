package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/repository"
	"github.com/redis/go-redis/v9"
)

type MovieService interface {
	AddMovie(movieDto *models.MovieDto) (*models.Movie, error)
	GetAllMovies() ([]models.MovieElastic, error)
	GetMovieByID(id string) (*models.MovieElastic, error)
	SearchMovies(query string) ([]models.MovieElastic, error)
	IngestMoviesFromDatabase()
}

type movieService struct {
	movieRepo        *repository.MovieRepository
	movieElasticRepo *repository.MovieElasticRepository
	redisClient      *redis.Client
}

func NewMovieService(movieRepo *repository.MovieRepository, movieElasticRepo *repository.MovieElasticRepository, redisClient *redis.Client) MovieService {
	return &movieService{
		movieRepo:        movieRepo,
		movieElasticRepo: movieElasticRepo,
		redisClient:      redisClient,
	}
}

func (s *movieService) AddMovie(movieDto *models.MovieDto) (*models.Movie, error) {
	movie := &models.Movie{
		Name:       movieDto.Name,
		Persona:    movieDto.Persona,
		MovieType:  movieDto.MovieType,
		LastUpdate: time.Now(),
	}

	err := s.movieRepo.Create(movie)
	if err != nil {
		return nil, err
	}

	// Also index in Elasticsearch
	movieElastic := &models.MovieElastic{
		ID:         movie.ID.Hex(),
		Name:       movie.Name,
		Persona:    movie.Persona,
		MovieType:  movie.MovieType,
		LastUpdate: movie.LastUpdate,
	}
	s.movieElasticRepo.Create(movieElastic)

	return movie, nil
}

func (s *movieService) GetAllMovies() ([]models.MovieElastic, error) {
	return s.movieElasticRepo.FindAll()
}

func (s *movieService) GetMovieByID(id string) (*models.MovieElastic, error) {
	return s.movieElasticRepo.FindByID(id)
}

func (s *movieService) SearchMovies(query string) ([]models.MovieElastic, error) {
	if query == "" {
		return s.movieElasticRepo.FindAll()
	}
	return s.movieElasticRepo.SearchByText(query)
}

func (s *movieService) IngestMoviesFromDatabase() {
	ctx := context.Background()

	// Get last updated time from Redis
	lastUpdatedTime := "1970-01-01T00:00:00Z"
	if val, err := s.redisClient.Get(ctx, "lastUpdatedTime").Result(); err == nil {
		lastUpdatedTime = val
	}

	// Find movies updated after lastUpdatedTime
	movies, err := s.movieRepo.FindAfterLastUpdatedTime(lastUpdatedTime)
	if err != nil {
		fmt.Printf("Error finding movies: %v\n", err)
		return
	}

	if len(movies) == 0 {
		return
	}

	// Convert to elasticsearch format
	var movieElastics []*models.MovieElastic
	for _, movie := range movies {
		movieElastics = append(movieElastics, &models.MovieElastic{
			ID:         movie.ID.Hex(),
			Name:       movie.Name,
			Persona:    movie.Persona,
			MovieType:  movie.MovieType,
			LastUpdate: movie.LastUpdate,
		})
	}

	// Save to Elasticsearch
	if err := s.movieElasticRepo.CreateBulk(movieElastics); err != nil {
		fmt.Printf("Error saving to Elasticsearch: %v\n", err)
		return
	}

	// Update last updated time in Redis
	currTime := time.Now().Format(time.RFC3339)
	s.redisClient.Set(ctx, "lastUpdatedTime", currTime, 0)

	fmt.Printf("Ingested %d movies at %s\n", len(movies), currTime)
}