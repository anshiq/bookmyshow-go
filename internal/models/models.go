package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Persona represents cast/crew members in a movie
type Persona struct {
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"` // lead, actor, producer, director
}

// SeatType represents a seat category with its price
type SeatType struct {
	Type  string `json:"type" bson:"type"`
	Price int    `json:"price" bson:"price"`
}

// Movie represents a movie in the system
type Movie struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Persona    []Persona          `json:"persona" bson:"persona"`
	MovieType  string             `json:"movieType" bson:"movieType"`
	LastUpdate time.Time          `json:"lastUpdated" bson:"lastUpdated"`
}

// MovieElastic represents a movie indexed in Elasticsearch
type MovieElastic struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Persona    []Persona `json:"persona"`
	MovieType  string    `json:"movieType"`
	LastUpdate time.Time `json:"lastUpdated"`
}

// Theatre represents a theatre/cinema hall
type Theatre struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Location string             `json:"location" bson:"location"`
	City     string             `json:"city" bson:"city"`
}

// Screen represents a screen within a theatre
type Screen struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	TheatreID string             `json:"theatreId" bson:"theatreId"`
}

// Show represents a movie show timing
type Show struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StartTime   time.Time         `json:"startTime" bson:"startTime"`
	EndTime     time.Time         `json:"endTime" bson:"endTime"`
	MovieID     string            `json:"movieId" bson:"movieId"`
	TheatreID   string            `json:"theatreId" bson:"theatreId"`
	ScreenID    string            `json:"screenId" bson:"screenId"`
	SeatType    []SeatType        `json:"seatType" bson:"seatType"`
	SeatMarking []int             `json:"seatMarking" bson:"seatMarking"`
}

// Ticket represents a booked ticket
type Ticket struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SeatNumber     string             `json:"seatNumber" bson:"seatNumber"`
	SeatCategory   string             `json:"seatCategory" bson:"seatCategory"`
	Price          int                `json:"price" bson:"price"`
	ShowID         string             `json:"showId" bson:"showId"`
	BookedByUserID string             `json:"bookedByUserId" bson:"bookedByUserId"`
	Status         bool               `json:"status" bson:"status"`
	PaymentStatus  bool               `json:"paymentStatus" bson:"paymentStatus"`
	PaymentMethod  string             `json:"paymentMethod" bson:"paymentMethod"`
}

// User represents a user of the system
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Age         int                `json:"age" bson:"age"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
}

// Config represents database configuration
type Config struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	HashID   string             `json:"hashId" bson:"hashId"`
	DBName   string             `json:"dbName" bson:"dbName"`
	DBType   string             `json:"dbType" bson:"dbType"`
	Username string             `json:"userName" bson:"userName"`
	Password string             `json:"password" bson:"password"`
	Host     string             `json:"host" bson:"host"`
}

// AllMovieShows represents combined movie-theatre-screen info
type AllMovieShows struct {
	MovieName       string    `json:"movieName"`
	Persona         []Persona `json:"persona"`
	MovieType       string    `json:"movieType"`
	TheatreName     string    `json:"theatreName"`
	TheatreLocation string    `json:"theatreLocation"`
	TheatreCity     string    `json:"theatreCity"`
	Screen          string    `json:"screen"`
}
