package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Seat represents a seat in a screen
type Seat struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SeatNumber    int                `json:"seatNumber" bson:"seatNumber"`
	SeatType      string             `json:"seatType" bson:"seatType"`
	ShowID        string             `json:"showId" bson:"showId"`
	BookingStatus bool               `json:"bookingStatus" bson:"bookingStatus"`
	Version       int                `json:"version" bson:"version"`
}
