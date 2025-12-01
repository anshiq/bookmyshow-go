package service

import (
	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/repository"
)

type TheatreService interface {
	AddTheatre(theatreDto *models.TheatreDto) (*models.Theatre, error)
	GetAllTheatres() ([]models.Theatre, error)
	GetTheatreByID(id string) (*models.Theatre, error)
}

type theatreService struct {
	theatreRepo *repository.TheatreRepository
}

func NewTheatreService(theatreRepo *repository.TheatreRepository) TheatreService {
	return &theatreService{theatreRepo: theatreRepo}
}

func (s *theatreService) AddTheatre(theatreDto *models.TheatreDto) (*models.Theatre, error) {
	theatre := &models.Theatre{
		Name:     theatreDto.Name,
		Location: theatreDto.Location,
		City:     theatreDto.City,
	}

	err := s.theatreRepo.Create(theatre)
	if err != nil {
		return nil, err
	}

	return theatre, nil
}

func (s *theatreService) GetAllTheatres() ([]models.Theatre, error) {
	return s.theatreRepo.FindAll()
}

func (s *theatreService) GetTheatreByID(id string) (*models.Theatre, error) {
	return s.theatreRepo.FindByID(id)
}