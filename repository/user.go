package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory(userID uint) ([]model.UserTaskCategory, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserTaskCategory(userID uint) ([]model.UserTaskCategory, error) {
	var userTaskCategories []model.UserTaskCategory
	if err := r.db.Table("users").
		Select("users.id, users.fullname, users.email, tasks.title as task, tasks.deadline, tasks.priority, tasks.status, categories.name as category").
		Joins("left join tasks on tasks.user_id = users.id").
		Joins("left join categories on tasks.category_id = categories.id").
		Where("users.id = ?", userID).
		Scan(&userTaskCategories).Error; err != nil {
		return nil, err
	}
	return userTaskCategories, nil
}

