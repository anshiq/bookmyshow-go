package service

import (
	"fmt"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/repository"
	"github.com/google/uuid"
)

type ShowService interface {
	AddShow(showDto *models.ShowDto) (*models.Show, error)
	GetAllShows() ([]models.Show, error)
	GetShowByID(id string) (*models.Show, error)
	GetShowsByMovieID(movieID string) ([]models.Show, error)
	GetMovieShows(movieID string) ([]models.AllMovieShows, error)
}

type showService struct {
	showRepo        *repository.ShowRepository
	movieElasticRepo *repository.MovieElasticRepository
	theatreRepo     *repository.TheatreRepository
	screenRepo      *repository.ScreenRepository
	seatRepo        *repository.SeatRepository
}

func NewShowService(showRepo *repository.ShowRepository, movieElasticRepo *repository.MovieElasticRepository, theatreRepo *repository.TheatreRepository, screenRepo *repository.ScreenRepository, seatRepo *repository.SeatRepository) ShowService {
	return &showService{
		showRepo:        showRepo,
		movieElasticRepo: movieElasticRepo,
		theatreRepo:     theatreRepo,
		screenRepo:      screenRepo,
		seatRepo:        seatRepo,
	}
}

func (s *showService) AddShow(showDto *models.ShowDto) (*models.Show, error) {
	// Validate movie exists in Elasticsearch
	_, err := s.movieElasticRepo.FindByID(showDto.MovieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Validate screen exists
	_, err = s.screenRepo.FindByID(showDto.ScreenID)
	if err != nil {
		return nil, fmt.Errorf("screen not found: %w", err)
	}

	// Validate theatre exists
	_, err = s.theatreRepo.FindByID(showDto.TheatreID)
	if err != nil {
		return nil, fmt.Errorf("theatre not found: %w", err)
	}

	show := &models.Show{
		StartTime:   showDto.StartTime,
		EndTime:     showDto.EndTime,
		MovieID:     showDto.MovieID,
		TheatreID:   showDto.TheatreID,
		ScreenID:    showDto.ScreenID,
		SeatType:    showDto.SeatType,
		SeatMarking: showDto.SeatMarking,
	}

	err = s.showRepo.Create(show)
	if err != nil {
		return nil, err
	}

	// Auto-generate seats based on seatMarking
	x := 0
	for i := 0; i < len(showDto.SeatMarking); i++ {
		for j := 0; j < showDto.SeatMarking[i]; j++ {
			x++
			seat := &models.Seat{
				SeatNumber:    x,
				SeatType:      showDto.SeatType[i].Type,
				ShowID:        show.ID.Hex(),
				BookingStatus: false,
				Version:       0,
			}
			s.seatRepo.Create(seat)
		}
	}

	return show, nil
}

func (s *showService) GetAllShows() ([]models.Show, error) {
	return s.showRepo.FindAll()
}

func (s *showService) GetShowByID(id string) (*models.Show, error) {
	return s.showRepo.FindByID(id)
}

func (s *showService) GetShowsByMovieID(movieID string) ([]models.Show, error) {
	return s.showRepo.FindByMovieID(movieID)
}

func (s *showService) GetMovieShows(movieID string) ([]models.AllMovieShows, error) {
	// Check if movie exists
	movie, err := s.movieElasticRepo.FindByID(movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found")
	}

	shows, err := s.showRepo.FindByMovieID(movieID)
	if err != nil {
		return nil, err
	}

	var allShows []models.AllMovieShows
	for _, show := range shows {
		theatre, err := s.theatreRepo.FindByID(show.TheatreID)
		if err != nil {
			continue
		}

		screen, err := s.screenRepo.FindByID(show.ScreenID)
		if err != nil {
			continue
		}

		allShows = append(allShows, models.AllMovieShows{
			MovieName:       movie.Name,
			Persona:         movie.Persona,
			MovieType:       movie.MovieType,
			TheatreName:     theatre.Name,
			TheatreLocation: theatre.Location,
			TheatreCity:     theatre.City,
			Screen:          screen.Name,
		})
	}

	return allShows, nil
}

var _ = uuid.New
