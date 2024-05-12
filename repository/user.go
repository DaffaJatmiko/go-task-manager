package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/db/filebased"
	"github.com/DaffaJatmiko/go-task-manager/model"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	filebasedDb *filebased.Data
}

func NewUserRepo(filebasedDb *filebased.Data) *userRepository {
	return &userRepository{filebasedDb}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	userByEmail, err := r.filebasedDb.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	return userByEmail, nil // TODO: replace this
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	createdUser, err := r.filebasedDb.CreateUser(user)

	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	userTask, err := r.filebasedDb.GetUserTaskCategory()
	if err != nil {
		return []model.UserTaskCategory{}, err
	}
	return userTask, nil // TODO: replace this
}
