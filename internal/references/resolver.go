package references

import (
	"io"
	"strings"

	"github.com/Sanssy/banter-engine/internal/logging"
	"github.com/Sanssy/banter-engine/internal/model"
)

const (
	clubPrefix = "mpp_championship_club_"
	userPrefix = "user_"
)

type Resolver struct {
	clubs  map[string]string
	users  map[string]string
	logger *logging.Logger
}

func New(output io.Writer) *Resolver {
	return &Resolver{
		clubs:  make(map[string]string),
		users:  make(map[string]string),
		logger: logging.New(output, "references"),
	}
}

func (r *Resolver) RegisterClub(id, name string) {
	if id != "" && name != "" {
		r.clubs[id] = name
	}
}

func (r *Resolver) RegisterUsers(standings []model.Standing) {
	for _, standing := range standings {
		r.RegisterUser(standing.UserID, standing.Name)
	}
	r.logger.Info("user reference loaded users_count=%d", len(r.users))
}

func (r *Resolver) RegisterUser(id, name string) {
	if id != "" && name != "" {
		r.users[id] = name
	}
}

func (r *Resolver) ClubName(id string) string {
	name, found := r.clubs[id]
	if found {
		return name
	}

	returnedName := id
	if strings.HasPrefix(id, clubPrefix) {
		returnedName = "Equipe inconnue"
	}
	return returnedName
}

func (r *Resolver) UserName(id string) string {
	name, found := r.users[id]
	if found {
		return name
	}
	if strings.HasPrefix(id, userPrefix) {
		return "Participant inconnu"
	}
	return id
}

func (r *Resolver) Resolve(value string) string {
	switch {
	case strings.HasPrefix(value, clubPrefix):
		return r.ClubName(value)
	case strings.HasPrefix(value, userPrefix):
		return r.UserName(value)
	default:
		return value
	}
}
