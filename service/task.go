package service

import (
	"context"

	"github.com/DaffaJatmiko/go-task-manager/config"
	"github.com/DaffaJatmiko/go-task-manager/model"
	repo "github.com/DaffaJatmiko/go-task-manager/repository"
)

type TaskService interface {
	Store(task *model.Task) error
	Update(id int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList() ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
}

type taskService struct {
	taskRepository repo.TaskRepository
}

func NewTaskService(taskRepository repo.TaskRepository) TaskService {
	return &taskService{taskRepository}
}

func (s *taskService) Store(task *model.Task) error {
	err := s.taskRepository.Store(task)
	if err == nil {
		config.RedisClient.Del(context.Background(), "taskList")
	}
	return err
}

func (s *taskService) Update(id int, task *model.Task) error {
	err := s.taskRepository.Update(id, task)
	if err == nil {
		config.RedisClient.Del(context.Background(), "taskList")
	}
	return err
}

func (s *taskService) Delete(id int) error {
	err := s.taskRepository.Delete(id)
	if err == nil {
		config.RedisClient.Del(context.Background(), "taskList")
	}
	return err
}

func (s *taskService) GetByID(id int) (*model.Task, error) {
	return s.taskRepository.GetByID(id)
}

func (s *taskService) GetList() ([]model.Task, error) {
	return s.taskRepository.GetList()
}

func (s *taskService) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	return s.taskRepository.GetTaskCategory(id)
}
