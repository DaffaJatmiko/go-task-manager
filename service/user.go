package service

import (
	"errors"
	"time"

	"github.com/DaffaJatmiko/go-task-manager/model"
	repo "github.com/DaffaJatmiko/go-task-manager/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *model.User) (model.User, error)
	Login(user *model.User) (token *string, err error)
	Logout(email string) error
	GetUserTaskCategory(userID uint) ([]model.UserTaskCategory, error)
}

type userService struct {
	userRepo     repo.UserRepository
	sessionsRepo repo.SessionRepository
}

func NewUserService(userRepository repo.UserRepository, sessionsRepo repo.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

func (s *userService) Register(user *model.User) (model.User, error) {
	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(user.Email)
	if err == nil {
		return model.User{}, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (s *userService) Login(user *model.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return nil, errors.New("wrong email or password")
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &model.Claims{
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(model.JwtKey)
	if err != nil {
		return nil, err
	}

	session := model.Session{
		Token:  tokenString,
		Email:  user.Email,
		Expiry: expirationTime,
	}

	existingSession, err := s.sessionsRepo.SessionAvailEmail(session.Email)
	if err != nil || existingSession.Email == "" {
		s.sessionsRepo.AddSessions(session)
	} else {
		s.sessionsRepo.UpdateSessions(session)
	}

	return &tokenString, nil
}

func (s *userService) Logout(email string) error {
	return s.sessionsRepo.DeleteSession(email)
}

func (s *userService) GetUserTaskCategory(userID uint) ([]model.UserTaskCategory, error) {
	taskByCategory, err := s.userRepo.GetUserTaskCategory(userID)
	if err != nil {
		return nil, err
	}
	return taskByCategory, nil
}
