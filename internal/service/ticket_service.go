package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/anshiq/bookmyshow-go/internal/models"
	"github.com/anshiq/bookmyshow-go/internal/payment"
	"github.com/anshiq/bookmyshow-go/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

type TicketService interface {
	BookTicket(ticketDto *models.TicketDto, userID string) (*models.Ticket, error)
	GetAllTickets() ([]models.TicketResponseDto, error)
	GetTicketsByUserID(userID string) ([]models.Ticket, error)
	CancelTicket(ticketID string, userID string) (string, error)
}

type ticketService struct {
	ticketRepo     *repository.TicketRepository
	seatRepo      *repository.SeatRepository
	showRepo      *repository.ShowRepository
	theatreRepo   *repository.TheatreRepository
	screenRepo    *repository.ScreenRepository
	movieElasticRepo *repository.MovieElasticRepository
	redisClient   *redis.Client
	paymentService *payment.PaymentService
}

func NewTicketService(ticketRepo *repository.TicketRepository, seatRepo *repository.SeatRepository, showRepo *repository.ShowRepository, theatreRepo *repository.TheatreRepository, screenRepo *repository.ScreenRepository, movieElasticRepo *repository.MovieElasticRepository, redisClient *redis.Client, paymentService *payment.PaymentService) TicketService {
	return &ticketService{
		ticketRepo:     ticketRepo,
		seatRepo:      seatRepo,
		showRepo:      showRepo,
		theatreRepo:   theatreRepo,
		screenRepo:    screenRepo,
		movieElasticRepo: movieElasticRepo,
		redisClient:   redisClient,
		paymentService: paymentService,
	}
}

func (s *ticketService) BookTicket(ticketDto *models.TicketDto, userID string) (*models.Ticket, error) {
	ctx := context.Background()

	// Parse seat number
	seatNumber := 0
	fmt.Sscanf(ticketDto.SeatNumber, "%d", &seatNumber)

	// Check if seat is available
	available, err := s.seatRepo.FindByNumberAndShowID(ticketDto.ShowID, seatNumber)
	if err != nil {
		return nil, errors.New("seat not found")
	}
	if available.BookingStatus {
		return nil, errors.New("seat not available")
	}

	// Get show details
	show, err := s.showRepo.FindByID(ticketDto.ShowID)
	if err != nil {
		return nil, err
	}

	// Find seat to get seat type
	seat, err := s.seatRepo.FindByNumberAndShowID(ticketDto.ShowID, seatNumber)
	if err != nil {
		return nil, err
	}

	// Calculate price based on seat type
	price := 250
	for _, st := range show.SeatType {
		if st.Type == seat.SeatType {
			price = st.Price
			break
		}
	}

	// Process payment
	paymentFactory := payment.NewPaymentProviderFactory()
	paymentService := paymentFactory.GetPaymentService()
	result, err := paymentService.ProcessPayment(ticketDto.PaymentMethod, float64(price), "INR", nil)
	if err != nil || result.Status != "success" {
		return nil, errors.New("payment failed")
	}

	// Change seat status
	s.seatRepo.Update(&models.Seat{
		SeatNumber:    seatNumber,
		BookingStatus: true,
	})

	ticket := &models.Ticket{
		SeatNumber:    ticketDto.SeatNumber,
		SeatCategory:  seat.SeatType,
		Price:         price,
		ShowID:        ticketDto.ShowID,
		BookedByUserID: userID,
		Status:        true,
		PaymentStatus: true,
		PaymentMethod: ticketDto.PaymentMethod,
	}

	err = s.ticketRepo.Create(ticket)
	if err != nil {
		return nil, err
	}

	// Mark seat as booked in Redis
	seatKey := fmt.Sprintf("seat:%s:%s", ticketDto.ShowID, ticketDto.SeatNumber)
	s.redisClient.Set(ctx, seatKey, ticket.ID.Hex(), 0)

	return ticket, nil
}

func (s *ticketService) GetAllTickets() ([]models.TicketResponseDto, error) {
	tickets, err := s.ticketRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []models.TicketResponseDto
	for _, ticket := range tickets {
		show, err := s.showRepo.FindByID(ticket.ShowID)
		if err != nil {
			continue
		}

		movie, err := s.movieElasticRepo.FindByID(show.MovieID)
		if err != nil {
			continue
		}

		theatre, err := s.theatreRepo.FindByID(show.TheatreID)
		if err != nil {
			continue
		}

		screen, err := s.screenRepo.FindByID(show.ScreenID)
		if err != nil {
			continue
		}

		response = append(response, models.TicketResponseDto{
			SeatNumber:   ticket.SeatNumber,
			SeatCategory: ticket.SeatCategory,
			Price:        ticket.Price,
			ShowDate:     show.StartTime,
			Description: models.AllMovieShows{
				MovieName:       movie.Name,
				Persona:         movie.Persona,
				MovieType:       movie.MovieType,
				TheatreName:     theatre.Name,
				TheatreLocation: theatre.Location,
				TheatreCity:     theatre.City,
				Screen:          screen.Name,
			},
			BookedStatus:  ticket.Status,
			PaymentStatus: ticket.PaymentStatus,
		})
	}

	return response, nil
}

func (s *ticketService) GetTicketsByUserID(userID string) ([]models.Ticket, error) {
	return s.ticketRepo.FindByUserID(userID)
}

func (s *ticketService) CancelTicket(ticketID string, userID string) (string, error) {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return "No Ticket Found", err
	}

	// Verify ownership
	if ticket.BookedByUserID != userID {
		return "Unauthorized", errors.New("unauthorized")
	}

	// Parse seat number
	seatNumber := 0
	fmt.Sscanf(ticket.SeatNumber, "%d", &seatNumber)

	// Release seat
	s.seatRepo.Update(&models.Seat{
		SeatNumber:    seatNumber,
		BookingStatus: false,
	})

	// Delete ticket
	s.ticketRepo.Delete(ticketID)

	// Release Redis lock
	ctx := context.Background()
	seatKey := fmt.Sprintf("seat:%s:%s", ticket.ShowID, ticket.SeatNumber)
	s.redisClient.Del(ctx, seatKey)

	return "Ticket Cancelled Successfully", nil
}

var _ = uuid.New
