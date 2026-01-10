package service

import (
	"errors"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatService interface {
	IfSeatAvailable(showID string, seatNumber int) (bool, error)
	FindByNumberAndShowID(showID string, seatNumber int) (*models.Seat, error)
	ChangeStatusOfSeat(bool bool, showID string, seatNumber int) (*models.Seat, error)
}

type seatService struct {
	seatRepo *repository.SeatRepository
}

func NewSeatService(seatRepo *repository.SeatRepository) SeatService {
	return &seatService{seatRepo: seatRepo}
}

func (s *seatService) IfSeatAvailable(showID string, seatNumber int) (bool, error) {
	seat, err := s.seatRepo.FindByNumberAndShowID(showID, seatNumber)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return seat.BookingStatus, nil
}

func (s *seatService) FindByNumberAndShowID(showID string, seatNumber int) (*models.Seat, error) {
	return s.seatRepo.FindByNumberAndShowID(showID, seatNumber)
}

func (s *seatService) ChangeStatusOfSeat(book bool, showID string, seatNumber int) (*models.Seat, error) {
	seat, err := s.seatRepo.FindByNumberAndShowID(showID, seatNumber)
	if err != nil {
		return nil, err
	}

	if seat != nil {
		if book {
			seat.BookingStatus = true
		} else {
			seat.BookingStatus = false
		}
		s.seatRepo.Update(seat)
	}

	savedSeat := seat
	version := seat.Version

	for i := 0; i < 3; i++ {
		updatedSeat, err := s.seatRepo.UpdateVersion(showID, seatNumber, version)
		if err != nil {
			if i == 2 {
				return nil, errors.New("seat cannot be booked")
			}
			continue
		}
		if updatedSeat != nil {
			savedSeat = updatedSeat
			break
		}
		if i == 2 {
			return nil, errors.New("seat cannot be booked")
		}
	}

	return savedSeat, nil
}
