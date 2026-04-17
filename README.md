# BookMyShow-go

A Go-based movie ticket booking system API inspired by BookMyShow, built with the Gin web framework.

## Features

- **Movie Management**: Add and search movies with Elasticsearch-powered search
- **Theatre & Screen Management**: Manage multiple theatres and screens
- **Show Scheduling**: Create and manage movie show timings
- **Ticket Booking**: Book tickets for shows with seat selection
- **User Authentication**: Signup and login for users
- **Payment Integration**: Support for Paytm and PhonePe payment providers
- **Caching**: Redis-based caching for improved performance

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Primary Database**: MongoDB
- **Search Engine**: Elasticsearch
- **Cache**: Redis
- **Authentication**: JWT (via middleware)

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/              # Configuration and database connections
│   ├── handler/             # HTTP request handlers
│   ├── middleware/          # HTTP middleware (CORS, auth, etc.)
│   ├── models/              # Data models and DTOs
│   ├── payment/             # Payment provider integration
│   ├── repository/          # Database access layer
│   └── service/             # Business logic layer
├── go.mod
└── README.md
```

## API Endpoints

### Admin Routes (`/admin`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/admin/movie` | Add a new movie |
| POST | `/admin/theatre` | Add a new theatre |
| POST | `/admin/screen` | Add a new screen |
| POST | `/admin/show` | Create a new show |

### User Routes (`/user`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/user/signup` | User registration |
| POST | `/user/login` | User login |
| GET | `/user/movies` | Get all movies |
| GET | `/user/movie` | Search movies |
| GET | `/user/shows` | Get movie shows |
| POST | `/user/book-tickets` | Book tickets |
| GET | `/user/tickets` | Get all tickets |
| DELETE | `/user/cancel-tickets` | Cancel a ticket |

### Config Routes (`/config`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/config/db-type` | Get config by database type |
| GET | `/config/get-by-type-and-hashId` | Get config by type and hash ID |

## Configuration

Environment variables required:

| Variable | Description |
|----------|-------------|
| `MONGODB_URI` | MongoDB connection URI |
| `REDIS_ADDR` | Redis server address |
| `ELASTICSEARCH_URL` | Elasticsearch server URL |
| `SERVER_PORT` | Server port (default: 8080) |
| `DATABASE_NAME` | Primary MongoDB database name |

## Getting Started

### Prerequisites

- Go 1.21+
- MongoDB
- Redis
- Elasticsearch

### Installation

1. Clone the repository:
```bash
git clone https://github.com/anshiq/bookmyshow-go.git
cd bookmyshow-go
```

2. Install dependencies:
```bash
go mod download
```

3. Set environment variables:
```bash
export MONGODB_URI="mongodb://localhost:27017"
export REDIS_ADDR="localhost:6379"
export ELASTICSEARCH_URL="http://localhost:9200"
export SERVER_PORT="8080"
export DATABASE_NAME="bookmyshow"
```

4. Run the server:
```bash
go run cmd/main.go
```

## Models

- **Movie**: id, name, persona, movieType, lastUpdated
- **Theatre**: id, name, location, city
- **Screen**: id, name, theatreId
- **Show**: id, startTime, endTime, movieId, theatreId, screenId, seatType, seatMarking
- **Ticket**: id, seatNumber, seatCategory, price, showId, bookedByUserId, status, paymentStatus, paymentMethod
- **User**: id, name, age, phoneNumber, email, password

## Payment Providers

The system supports multiple payment providers via a factory pattern:
- Paytm
- PhonePe

## License

MIT
