package repository

import (
	"github.com/DaffaJatmiko/go-task-manager/db/filebased"
	"github.com/DaffaJatmiko/go-task-manager/model"
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
	filebased *filebased.Data
}

func NewTaskRepo(filebasedDb *filebased.Data) *taskRepository {
	return &taskRepository{
		filebased: filebasedDb,
	}
}

func (t *taskRepository) Store(task *model.Task) error {
	t.filebased.StoreTask(*task)

	return nil
}

func (t *taskRepository) Update(taskID int, task *model.Task) error {
	err := t.filebased.UpdateTask(taskID, *task)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (t *taskRepository) Delete(id int) error {
	err := t.filebased.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	task, err := t.filebased.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil // TODO: replace this
}

func (t *taskRepository) GetList() ([]model.Task, error) {
	taskList, err := t.filebased.GetTasks()
	if err != nil {
		return nil, err
	}
	return taskList, nil // TODO: replace this
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	taskByCategory, err := t.filebased.GetTaskListByCategory(id)
	if err != nil {
		return nil, err
	}
	return taskByCategory, nil // TODO: replace this
}
