# Iteration 001 - Fetch MPP Standings

## Objective
Implement the first working version of Banter Engine.

- Read MPP_TOKEN from environment variables
- Call the MPP API
- Retrieve standings for challenge mpp_challenge_UDKDDH27
- Map API data to domain models
- Print standings in the terminal

## Constraints

Do NOT introduce:
- Hexagonal Architecture
- Repository Pattern
- Service Layer
- CQRS
- Event Bus
- DI Frameworks
- Unnecessary interfaces

Prefer:
- Simple code
- Explicit code
- Small functions
- Standard library only

## Project Structure

cmd/
└── banter-engine/
    └── main.go

internal/
├── model/
│   └── standing.go
└── mpp/
    └── client.go

## Domain Model

```go
type Standing struct {
    UserID string
    Name   string
    Rank   int
    Points int
}
```

## API

Base URL:
https://api.mpp.football

Endpoint:
GET /challenge-standings/users-standings

Authentication:
Authorization: Bearer <MPP_TOKEN>

## Client API

```go
func (c *Client) GetStandings(challengeID string) ([]model.Standing, error)
```

Responsibilities:
- Perform HTTP request
- Deserialize JSON
- Map DTOs to domain models
- Return standings

## Definition of Done

The following command works:

go run ./cmd/banter-engine

And prints the standings for:
mpp_challenge_UDKDDH27
