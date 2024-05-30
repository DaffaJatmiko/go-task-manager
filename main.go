package main

import (
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"embed"
	"fmt"
	"net/http"

	"github.com/DaffaJatmiko/go-task-manager/client"
	"github.com/DaffaJatmiko/go-task-manager/config"
	"github.com/DaffaJatmiko/go-task-manager/handler/api"
	"github.com/DaffaJatmiko/go-task-manager/handler/web"
	"github.com/DaffaJatmiko/go-task-manager/middleware"
	"github.com/DaffaJatmiko/go-task-manager/migrations" // Import migrations package
	repo "github.com/DaffaJatmiko/go-task-manager/repository"
	"github.com/DaffaJatmiko/go-task-manager/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type APIHandler struct {
	UserAPIHandler     api.UserAPI
	CategoryAPIHandler api.CategoryAPI
	TaskAPIHandler     api.TaskAPI
}

type ClientHandler struct {
	AuthWeb      web.AuthWeb
	HomeWeb      web.HomeWeb
	DashboardWeb web.DashboardWeb
	TaskWeb      web.TaskWeb
	CategoryWeb  web.CategoryWeb
	ModalWeb     web.ModalWeb
}

//go:embed views/*
var Resources embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()
	config.InitRedis() 

	// Jalankan migrasi
	if err := migrations.Migrate(config.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	gin.SetMode(gin.ReleaseMode) //release

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := gin.New()
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] \"%s %s %s\"\n",
				param.TimeStamp.Format(time.RFC822),
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		}))
		router.Use(gin.Recovery())

		router = RunServer(router, config.DB)
		router = RunClient(router, Resources, config.DB)

		fmt.Println("Server is running on port 8080")
		err = router.Run(":8080")
		if err != nil {
			panic(err)
		}

	}()

	wg.Wait()
}

func RunServer(gin *gin.Engine, db *gorm.DB) *gin.Engine {
	userRepo := repo.NewUserRepo(db)
	sessionRepo := repo.NewSessionsRepo(db)
	categoryRepo := repo.NewCategoryRepo(db)
	taskRepo := repo.NewTaskRepo(db)

	userService := service.NewUserService(userRepo, sessionRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	taskService := service.NewTaskService(taskRepo)

	userAPIHandler := api.NewUserAPI(userService)
	categoryAPIHandler := api.NewCategoryAPI(categoryService)
	taskAPIHandler := api.NewTaskAPI(taskService)

	apiHandler := APIHandler{
		UserAPIHandler:     userAPIHandler,
		CategoryAPIHandler: categoryAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
	}

	version := gin.Group("/api/v1")
	{
		user := version.Group("/user")
		{
			user.POST("/login", apiHandler.UserAPIHandler.Login)
			user.POST("/register", apiHandler.UserAPIHandler.Register)

			user.Use(middleware.Auth())
			user.GET("/tasks/:id", apiHandler.UserAPIHandler.GetUserTaskCategory)
			user.POST("/logout", apiHandler.UserAPIHandler.Logout)
		}

		task := version.Group("/task")
		{
			task.Use(middleware.Auth())
			task.POST("/add", apiHandler.TaskAPIHandler.AddTask)
			task.GET("/get/:id", apiHandler.TaskAPIHandler.GetTaskByID)
			task.PUT("/update/:id", apiHandler.TaskAPIHandler.UpdateTask)
			task.DELETE("/delete/:id", apiHandler.TaskAPIHandler.DeleteTask)
			task.GET("/list", apiHandler.TaskAPIHandler.GetTaskList)
			task.GET("/category/:id", apiHandler.TaskAPIHandler.GetTaskListByCategory)
		}

		category := version.Group("/category")
		{
			category.Use(middleware.Auth())
			category.POST("/add", apiHandler.CategoryAPIHandler.AddCategory)
			category.GET("/get/:id", apiHandler.CategoryAPIHandler.GetCategoryByID)
			category.PUT("/update/:id", apiHandler.CategoryAPIHandler.UpdateCategory)
			category.DELETE("/delete/:id", apiHandler.CategoryAPIHandler.DeleteCategory)
			category.GET("/list", apiHandler.CategoryAPIHandler.GetCategoryList)
		}
	}

	return gin
}

func RunClient(gin *gin.Engine, embed embed.FS, db *gorm.DB) *gin.Engine {
	sessionRepo := repo.NewSessionsRepo(db)
	sessionService := service.NewSessionService(sessionRepo)

	userClient := client.NewUserClient()
	taskClient := client.NewTaskClient()
	categoryClient := client.NewCategoryClient()

	authWeb := web.NewAuthWeb(userClient, sessionService, embed)
	modalWeb := web.NewModalWeb(embed)
	homeWeb := web.NewHomeWeb(embed)
	dashboardWeb := web.NewDashboardWeb(userClient, sessionService, embed)
	taskWeb := web.NewTaskWeb(taskClient, sessionService, embed)
	categoryWeb := web.NewCategoryWeb(categoryClient, sessionService, embed)

	client := ClientHandler{
		authWeb, homeWeb, dashboardWeb, taskWeb, categoryWeb, modalWeb,
	}

	gin.StaticFS("/static", http.Dir("frontend/public"))

	gin.GET("/", client.HomeWeb.Index)

	user := gin.Group("/client")
	{
		user.GET("/login", client.AuthWeb.Login)
		user.POST("/login/process", client.AuthWeb.LoginProcess)
		user.GET("/register", client.AuthWeb.Register)
		user.POST("/register/process", client.AuthWeb.RegisterProcess)

		user.Use(middleware.Auth())
		user.GET("/logout", client.AuthWeb.Logout)
	}

	main := gin.Group("/client")
	{
		main.Use(middleware.Auth())
		main.GET("/dashboard", client.DashboardWeb.Dashboard)
		main.GET("/task", client.TaskWeb.TaskPage)
		main.POST("/task/add/process", client.TaskWeb.TaskAddProcess)
		main.GET("/category", client.CategoryWeb.Category)
	}

	modal := gin.Group("/client")
	{
		modal.GET("/modal", client.ModalWeb.Modal)
	}

	return gin
}
