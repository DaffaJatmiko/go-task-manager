package api

import (
	"github.com/DaffaJatmiko/go-task-manager/model"
	"github.com/DaffaJatmiko/go-task-manager/service"
	"net/http"
	"strconv"

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
	categoryList, err := ct.categoryService.GetList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, categoryList)
}
