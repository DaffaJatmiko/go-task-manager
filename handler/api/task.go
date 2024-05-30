package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DaffaJatmiko/go-task-manager/config"
	"github.com/DaffaJatmiko/go-task-manager/model"
	"github.com/DaffaJatmiko/go-task-manager/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type TaskAPI interface {
	AddTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetTaskByID(c *gin.Context)
	GetTaskList(c *gin.Context)
	GetTaskListByCategory(c *gin.Context)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) AddTask(c *gin.Context) {
	var newTask model.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	err := t.taskService.Store(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	// Hapus cache setelah menambahkan task baru
	config.RedisClient.Del(context.Background(), "taskList")

	c.JSON(http.StatusCreated, model.NewSuccessResponse("add task success"))
}

func (t *taskAPI) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid task ID"))
		return
	}

	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	task.ID = taskID
	if err = t.taskService.Update(task.ID, &task); err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	// Hapus cache setelah mengupdate task
	config.RedisClient.Del(context.Background(), "taskList")

	c.JSON(http.StatusOK, model.NewSuccessResponse("update task success"))
}

func (t *taskAPI) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Invalid task ID"))
		return
	}

	if err := t.taskService.Delete(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	// Hapus cache setelah menghapus task
	config.RedisClient.Del(context.Background(), "taskList")

	c.JSON(http.StatusOK, model.NewSuccessResponse("Task deleted successfully"))
}

func (t *taskAPI) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Invalid task ID"))
		return
	}

	task, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, task)
}

func (t *taskAPI) GetTaskList(c *gin.Context) {
	ctx := context.Background()

	// Cek cache terlebih dahulu
	cachedTasks, err := config.RedisClient.Get(ctx, "taskList").Result()
	if err == redis.Nil {
		// Cache tidak ditemukan, ambil dari database
		taskList, err := t.taskService.GetList()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
			return
		}

		// Simpan hasil ke cache
		taskData, err := json.Marshal(taskList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
			return
		}
		config.RedisClient.Set(ctx, "taskList", taskData, 5*time.Minute)

		log.Println("Task list fetched from database")
		c.JSON(http.StatusOK, taskList)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	} else {
		// Cache ditemukan, kembalikan dari cache
		var tasks []model.Task
		err := json.Unmarshal([]byte(cachedTasks), &tasks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
			return
		}

		log.Println("Task list fetched from cache")
		c.JSON(http.StatusOK, tasks)
	}
}

func (t *taskAPI) GetTaskListByCategory(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid category ID"))
		return
	}

	taskList, err := t.taskService.GetTaskCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, taskList)
}
