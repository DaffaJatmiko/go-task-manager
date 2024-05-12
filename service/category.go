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

func (c *categoryService) Store(category *model.Category) error {
	err := c.categoryRepository.Store(category)
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryService) Update(id int, category model.Category) error {
	err := c.categoryRepository.Update(id, category)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (c *categoryService) Delete(id int) error {
	err := c.categoryRepository.Delete(id)
	if err != nil {
		return err 
	}
	return nil // TODO: replace this
}

func (c *categoryService) GetByID(id int) (*model.Category, error) {
	category, err := c.categoryRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *categoryService) GetList() ([]model.Category, error) {
	categoryList, err := c.categoryRepository.GetList()
	if err != nil {
		return nil, err
	}
	return categoryList, nil // TODO: replace this
}
