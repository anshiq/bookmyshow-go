package service

import (
	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/repository"
)

type ScreenService interface {
	AddScreen(screenDto *models.ScreenDto) (*models.Screen, error)
	GetAllScreens() ([]models.Screen, error)
	GetScreenByID(id string) (*models.Screen, error)
	GetScreensByTheatreID(theatreID string) ([]models.Screen, error)
}

type screenService struct {
	screenRepo *repository.ScreenRepository
}

func NewScreenService(screenRepo *repository.ScreenRepository) ScreenService {
	return &screenService{screenRepo: screenRepo}
}

func (s *screenService) AddScreen(screenDto *models.ScreenDto) (*models.Screen, error) {
	screen := &models.Screen{
		Name:      screenDto.Name,
		TheatreID: screenDto.TheatreID,
	}

	err := s.screenRepo.Create(screen)
	if err != nil {
		return nil, err
	}

	return screen, nil
}

func (s *screenService) GetAllScreens() ([]models.Screen, error) {
	return s.screenRepo.FindAll()
}

func (s *screenService) GetScreenByID(id string) (*models.Screen, error) {
	return s.screenRepo.FindByID(id)
}

func (s *screenService) GetScreensByTheatreID(theatreID string) ([]models.Screen, error) {
	return s.screenRepo.FindByTheatreID(theatreID)
}