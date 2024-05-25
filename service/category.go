package service

import (
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
	return s.categoryRepository.Store(category)
}

func (s *categoryService) Update(id int, category model.Category) error {
	return s.categoryRepository.Update(id, category)
}

func (s *categoryService) Delete(id int) error {
	return s.categoryRepository.Delete(id)
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	return s.categoryRepository.GetByID(id)
}

func (s *categoryService) GetList() ([]model.Category, error) {
	return s.categoryRepository.GetList()
}
