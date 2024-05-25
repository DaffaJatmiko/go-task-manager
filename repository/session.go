package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	"gorm.io/gorm"
	"time"
)

type SessionRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(email string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) SessionRepository {
	return &sessionsRepo{db}
}

func (u *sessionsRepo) AddSessions(session model.Session) error {
	return u.db.Create(&session).Error
}

func (u *sessionsRepo) DeleteSession(email string) error {
	return u.db.Where("email = ?", email).Delete(&model.Session{}).Error
}

func (u *sessionsRepo) UpdateSessions(session model.Session) error {
	return u.db.Save(&session).Error
}

func (u *sessionsRepo) SessionAvailEmail(email string) (model.Session, error) {
	var session model.Session
	if err := u.db.Where("email = ?", email).First(&session).Error; err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (u *sessionsRepo) SessionAvailToken(token string) (model.Session, error) {
	var session model.Session
	if err := u.db.Where("token = ?", token).First(&session).Error; err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (u *sessionsRepo) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		if err := u.DeleteSession(token); err != nil {
			return model.Session{}, err
		}
		return model.Session{}, err
	}

	return session, nil
}

func (u *sessionsRepo) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
