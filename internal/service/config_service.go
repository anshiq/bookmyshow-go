package service

import (
	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/repository"
)

type ConfigService interface {
	GetByDBType(dbType string) ([]models.Config, error)
	GetByDBTypeAndHashID(dbType, hashID string) (*models.Config, error)
	SetConfig(configDto *models.ConfigDto) (*models.Config, error)
}

type configService struct {
	configRepo *repository.ConfigRepository
}

func NewConfigService(configRepo *repository.ConfigRepository) ConfigService {
	return &configService{configRepo: configRepo}
}

func (s *configService) GetByDBType(dbType string) ([]models.Config, error) {
	return s.configRepo.FindByDBType(dbType)
}

func (s *configService) GetByDBTypeAndHashID(dbType, hashID string) (*models.Config, error) {
	return s.configRepo.FindByDBTypeAndHashID(dbType, hashID)
}

func (s *configService) SetConfig(configDto *models.ConfigDto) (*models.Config, error) {
	config := &models.Config{
		HashID:   configDto.HashID,
		DBName:   configDto.DBName,
		DBType:   configDto.DBType,
		Username: configDto.Username,
		Password: configDto.Password,
		Host:     configDto.Host,
	}

	err := s.configRepo.Create(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
