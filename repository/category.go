package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Store(category *model.Category) error
	Update(id int, category model.Category) error
	Delete(id int) error
	GetByID(id int) (*model.Category, error)
	GetList() ([]model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (c *categoryRepository) Store(category *model.Category) error {
	return c.db.Create(category).Error
}

func (c *categoryRepository) Update(id int, category model.Category) error {
	return c.db.Model(&model.Category{}).Where("id = ?", id).Updates(category).Error
}

func (c *categoryRepository) Delete(id int) error {
	return c.db.Delete(&model.Category{}, id).Error
}

func (c *categoryRepository) GetByID(id int) (*model.Category, error) {
	var category model.Category
	if err := c.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *categoryRepository) GetList() ([]model.Category, error) {
	var categories []model.Category
	if err := c.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
