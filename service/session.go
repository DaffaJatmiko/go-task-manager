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

func NewSessionService(sessionRepo repo.SessionRepository) SessionService {
	return &sessionService{sessionRepo}
}

func (s *sessionService) GetSessionByEmail(email string) (model.Session, error) {
	session, err := s.sessionRepo.SessionAvailEmail(email)
	if err != nil {
		return model.Session{}, err
	}
	return session, nil
}
