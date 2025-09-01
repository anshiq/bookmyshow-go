package models

import "net/http"

type ErrorResults int

const (
	MOVIE_NOT_FOUND ErrorResults = iota
	THEATRE_NOT_FOUND
	SCREEN_NOT_FOUND
	SHOW_NOT_FOUND
	TICKET_NOT_FOUND
	USER_NOT_FOUND
	SHOWtime_NOT_AVAILABLE
	SEAT_ALREADY_BOOKED
	PAYMENT_FAILED
	INVALID_REQUEST
	USER_ALREADY_EXISTS
	INVALID_CREDENTIALS
	UNAUTHORIZED
)

func (e ErrorResults) Status() int {
	switch e {
	case MOVIE_NOT_FOUND, THEATRE_NOT_FOUND, SCREEN_NOT_FOUND, SHOW_NOT_FOUND, TICKET_NOT_FOUND, USER_NOT_FOUND:
		return http.StatusNotFound
	case SHOWtime_NOT_AVAILABLE, SEAT_ALREADY_BOOKED:
		return http.StatusConflict
	case PAYMENT_FAILED:
		return http.StatusPaymentRequired
	case INVALID_REQUEST, INVALID_CREDENTIALS:
		return http.StatusBadRequest
	case USER_ALREADY_EXISTS:
		return http.StatusConflict
	case UNAUTHORIZED:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (e ErrorResults) Message() string {
	switch e {
	case MOVIE_NOT_FOUND:
		return "Movie not found"
	case THEATRE_NOT_FOUND:
		return "Theatre not found"
	case SCREEN_NOT_FOUND:
		return "Screen not found"
	case SHOW_NOT_FOUND:
		return "Show not found"
	case TICKET_NOT_FOUND:
		return "Ticket not found"
	case USER_NOT_FOUND:
		return "User not found"
	case SHOWtime_NOT_AVAILABLE:
		return "Show time is not available"
	case SEAT_ALREADY_BOOKED:
		return "Seat is already booked"
	case PAYMENT_FAILED:
		return "Payment failed"
	case INVALID_REQUEST:
		return "Invalid request"
	case USER_ALREADY_EXISTS:
		return "User already exists"
	case INVALID_CREDENTIALS:
		return "Invalid credentials"
	case UNAUTHORIZED:
		return "Unauthorized"
	default:
		return "Internal server error"
	}
}

type PostExceptions struct {
	ErrorResults ErrorResults
}

func (e PostExceptions) Error() string {
	return e.ErrorResults.Message()
}

func NewPostException(err ErrorResults) *PostExceptions {
	return &PostExceptions{ErrorResults: err}
}
