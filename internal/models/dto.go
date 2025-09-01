package models

import "time"

// MovieDto - Data Transfer Object for creating a movie
type MovieDto struct {
	Name       string     `json:"name" binding:"required"`
	Persona    []Persona  `json:"persona"`
	MovieType  string     `json:"movieType"`
	LastUpdate *time.Time `json:"lastUpdated"`
}

// TheatreDto - Data Transfer Object for creating a theatre
type TheatreDto struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
	City     string `json:"city" binding:"required"`
}

// ScreenDto - Data Transfer Object for creating a screen
type ScreenDto struct {
	Name      string `json:"name" binding:"required"`
	TheatreID string `json:"theatreId" binding:"required"`
}

// ShowDto - Data Transfer Object for creating a show
type ShowDto struct {
	StartTime   time.Time  `json:"startTime" binding:"required"`
	EndTime     time.Time  `json:"endTime" binding:"required"`
	MovieID     string     `json:"movieId" binding:"required"`
	TheatreID   string     `json:"theatreId" binding:"required"`
	ScreenID    string     `json:"screenId" binding:"required"`
	SeatType    []SeatType `json:"seatType"`
	SeatMarking []int      `json:"seatMarking"`
}

// TicketDto - Data Transfer Object for booking a ticket
type TicketDto struct {
	ShowID        string `json:"showId" binding:"required"`
	SeatNumber    string `json:"seatNumber" binding:"required"`
	PaymentMethod string `json:"paymentMethod"`
}

// TicketResponseDto - Response DTO for ticket booking
type TicketResponseDto struct {
	SeatNumber    string        `json:"seatNumber"`
	SeatCategory  string        `json:"seatCategory"`
	Price         int           `json:"price"`
	ShowDate      time.Time     `json:"showDate"`
	Description   AllMovieShows `json:"description"`
	BookedStatus  bool          `json:"bookedStatus"`
	PaymentStatus bool          `json:"paymentStatus"`
}

// UserDto - Data Transfer Object for user registration
type UserDto struct {
	Name        string `json:"name" binding:"required"`
	Age         int    `json:"age"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

// ConfigDto - Data Transfer Object for database configuration
type ConfigDto struct {
	DBName   string `json:"dbName" binding:"required"`
	HashID   string `json:"hashId" binding:"required"`
	DBType   string `json:"dbType" binding:"required"`
	Username string `json:"userName"`
	Password string `json:"password"`
	Host     string `json:"host" binding:"required"`
}

// LoginDto - Data Transfer Object for user login
type LoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ErrorResponse - Standard error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}