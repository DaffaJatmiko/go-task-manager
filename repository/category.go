package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/db/filebased"
	"github.com/DaffaJatmiko/go-task-manager/model"
)

type CategoryRepository interface {
	Store(Category *model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	GetByID(id int) (*model.Category, error)
	GetList() ([]model.Category, error)
}

type categoryRepository struct {
	filebasedDb *filebased.Data
}

func NewCategoryRepo(filebasedDb *filebased.Data) *categoryRepository {
	return &categoryRepository{filebasedDb}
}

func (c *categoryRepository) Store(Category *model.Category) error {
	c.filebasedDb.StoreCategory(*Category)
	return nil
}

func (c *categoryRepository) Update(id int, category model.Category) error {
	err := c.filebasedDb.UpdateCategory(id, category)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (c *categoryRepository) Delete(id int) error {
	err := c.filebasedDb.DeleteCategory(id)
	if err != nil {
		return err 
	}
	return nil // TODO: replace this
}

func (c *categoryRepository) GetByID(id int) (*model.Category, error) {
	category, err := c.filebasedDb.GetCategoryByID(id)

	return category, err
}

func (c *categoryRepository) GetList() ([]model.Category, error) {
	categoryList, err := c.filebasedDb.GetCategories()
	if err != nil {
		return nil, err
	}
	return categoryList, nil // TODO: replace this
}
