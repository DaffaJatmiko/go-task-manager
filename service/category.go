package service

import (
	"context"

	"github.com/DaffaJatmiko/go-task-manager/config"
	"github.com/DaffaJatmiko/go-task-manager/model"
	repo "github.com/DaffaJatmiko/go-task-manager/repository"
)

type CategoryService interface {
	Store(category *model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	GetByID(id int) (*model.Category, error)
	GetList() ([]model.Category, error)
}

type categoryService struct {
	categoryRepository repo.CategoryRepository
}

func NewCategoryService(categoryRepository repo.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository}
}

func (s *categoryService) Store(category *model.Category) error {
	err := s.categoryRepository.Store(category)
	if err == nil {
		config.RedisClient.Del(context.Background(), "categoryList")
	}
	return nil
}

func (s *categoryService) Update(id int, category model.Category) error {
	err := s.categoryRepository.Update(id, category)
	if err == nil {
		config.RedisClient.Del(context.Background(), "categoryList")
	}
	return nil
}

func (s *categoryService) Delete(id int) error {
	err := s.categoryRepository.Delete(id)
	if err == nil {
		config.RedisClient.Del(context.Background(), "categoryList")
	}
	return nil
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	return s.categoryRepository.GetByID(id)
}

func (s *categoryService) GetList() ([]model.Category, error) {
	return s.categoryRepository.GetList()
}
