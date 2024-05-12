package service

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	repo "github.com/DaffaJatmiko/go-task-manager/repository"
)

type SessionService interface {
	GetSessionByEmail(email string) (model.Session, error)
}

type sessionService struct {
	sessionRepo repo.SessionRepository
}

func NewSessionService(sessionRepo repo.SessionRepository) *sessionService {
	return &sessionService{sessionRepo}
}

func (c *sessionService) GetSessionByEmail(email string) (model.Session, error) {
	session, err := c.sessionRepo.SessionAvailEmail(email)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil // TODO: replace this
}
