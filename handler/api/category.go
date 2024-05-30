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
	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
)

type CategoryAPI interface {
	AddCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	GetCategoryList(c *gin.Context)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (ct *categoryAPI) AddCategory(c *gin.Context) {
	var newCategory model.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	err := ct.categoryService.Store(&newCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("add category success"))
}

func (ct *categoryAPI) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid Category ID"))
		return
	}

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(err.Error()))
		return
	}

	category.ID = categoryID
	if err := ct.categoryService.Update(category.ID, category); err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("category update success"))
}

func (ct *categoryAPI) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid Category ID"))
		return
	}

	if err := ct.categoryService.Delete(categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse("Category deleted successfully"))
}

func (ct *categoryAPI) GetCategoryByID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Invalid category ID"))
		return
	}

	category, err := ct.categoryService.GetByID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, category)
}

func (ct *categoryAPI) GetCategoryList(c *gin.Context) {
	ctx := context.Background()

	//cek cache terlebih dahulu
	cachedCategories, err := config.RedisClient.Get(ctx, "categoryList").Result()
	if err == redis.Nil {
		//cache tidak ditemukan, ambil dari database
		categoryList, err := ct.categoryService.GetList()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
			return
		}

		//simpan ke cache
		categoryData, err := json.Marshal(categoryList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
			return
		}
		config.RedisClient.Set(ctx, "categoryList", categoryData, 5*time.Minute)

		log.Println("Category list fetched from database")
		c.JSON(http.StatusOK, categoryList)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	} else {
		//cache ditemukan, kembalikan dari cache
		var categories []model.Category
		err = json.Unmarshal([]byte(cachedCategories), &categories)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
			return
		}
		log.Println("Category list fetched from cache")
		c.JSON(http.StatusOK, categories)
	}
}
