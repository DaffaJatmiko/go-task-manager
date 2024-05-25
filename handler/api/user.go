package api

import (
	"net/http"
	"strconv"

	"github.com/DaffaJatmiko/go-task-manager/model"
	"github.com/DaffaJatmiko/go-task-manager/service"

	"github.com/gin-gonic/gin"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	var user model.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("email or password is empty"))
		return
	}

	var recordUser = model.User{
		Email:    user.Email,
		Password: user.Password,
	}

	tokenStr, err := u.userService.Login(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.SetCookie("session_token", *tokenStr, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login success"})
}

func (u *userAPI) Logout(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
		return
	}

	err := u.userService.Logout(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("failed to logout"))
		return
	}

	// Clear the session cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("Invalid user ID"))
		return
	}

	userTask, err := u.userService.GetUserTaskCategory(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusOK, userTask)
}
