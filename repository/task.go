package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(taskID int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (t *taskRepository) Store(task *model.Task) error {
	return t.db.Create(task).Error
}

func (t *taskRepository) Update(taskID int, task *model.Task) error {
	return t.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(task).Error
}

func (t *taskRepository) Delete(id int) error {
	return t.db.Delete(&model.Task{}, id).Error
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	var task model.Task
	if err := t.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskRepository) GetList() ([]model.Task, error) {
	var tasks []model.Task
	if err := t.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	var taskCategories []model.TaskCategory
	if err := t.db.Table("tasks").
		Select("tasks.id, tasks.title, categories.name as category").
		Joins("left join categories on tasks.category_id = categories.id").
		Where("tasks.category_id = ?", id).
		Scan(&taskCategories).Error; err != nil {
		return nil, err
	}
	return taskCategories, nil
}
