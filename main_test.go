package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/DaffaJatmiko/go-task-manager/handler/api"
	"github.com/DaffaJatmiko/go-task-manager/model"
	"github.com/DaffaJatmiko/go-task-manager/repository"
	"github.com/DaffaJatmiko/go-task-manager/service"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router          *gin.Engine
	userAPI         api.UserAPI
	taskAPI         api.TaskAPI
	categoryAPI     api.CategoryAPI
	userService     service.UserService
	taskService     service.TaskService
	categoryService service.CategoryService
	userRepo        repository.UserRepository
	taskRepo        repository.TaskRepository
	categoryRepo    repository.CategoryRepository
	sessionRepo     repository.SessionRepository
	db              *gorm.DB
	token           string

)

func setupRouter() *gin.Engine {
	router := gin.Default()
	apiRoutes := router.Group("/api/v1")
	{
		apiRoutes.POST("/user/register", userAPI.Register)
		apiRoutes.POST("/user/login", userAPI.Login)
		apiRoutes.GET("/task/:id", userAPI.GetUserTaskCategory)
		apiRoutes.POST("/user/logout", userAPI.Logout)
		apiRoutes.POST("/tasks", taskAPI.AddTask)
		apiRoutes.GET("/tasks/:id", taskAPI.GetTaskByID)
		apiRoutes.GET("/tasks", taskAPI.GetTaskList)
		apiRoutes.PUT("/tasks/:id", taskAPI.UpdateTask)
		apiRoutes.DELETE("/tasks/:id", taskAPI.DeleteTask)
		apiRoutes.POST("/category", categoryAPI.AddCategory)
		apiRoutes.GET("/categories/:id", categoryAPI.GetCategoryByID)
		apiRoutes.GET("/categories", categoryAPI.GetCategoryList)
		apiRoutes.PUT("/categories/:id", categoryAPI.UpdateCategory)
		apiRoutes.DELETE("/categories/:id", categoryAPI.DeleteCategory)
	}
	return router
}

func CreateTestTask(id int, title, deadline string, priority int, status string, categoryID int, userID int) {
	task := model.Task{
			ID:         id,
			Title:      title,
			Deadline:   deadline,
			Priority:   priority,
			Status:     status,
			CategoryID: categoryID,
			UserID:     userID,
	}
	// Simpan task ke database
	db.Create(&task)
}

func CreateTestCategory(id int, name string) {
	category := model.Category{
			ID:   id,
			Name: name,
	}
	// Simpan category ke database
	db.Create(&category)
}


var _ = BeforeSuite(func() {
	// Set up the test database connection directly in the test file
	dsn := "host=localhost user=postgres password=jatming dbname=test-task-manager port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	Expect(err).NotTo(HaveOccurred())

	// Apply migrations to the test database
	err = db.AutoMigrate(&model.User{}, &model.Task{}, &model.Category{}, &model.Session{})
	if err != nil {
		log.Fatalf("failed to migrate test database: %v", err)
	}
	Expect(err).NotTo(HaveOccurred())

	// Set up the repositories, services, and APIs
	userRepo = repository.NewUserRepo(db)
	sessionRepo = repository.NewSessionsRepo(db)
	taskRepo = repository.NewTaskRepo(db)
	categoryRepo = repository.NewCategoryRepo(db)
	userService = service.NewUserService(userRepo, sessionRepo)
	taskService = service.NewTaskService(taskRepo)
	categoryService = service.NewCategoryService(categoryRepo)
	userAPI = api.NewUserAPI(userService)
	taskAPI = api.NewTaskAPI(taskService)
	categoryAPI = api.NewCategoryAPI(categoryService)

	// Set up the router
	router = setupRouter()

	
	// Ensure test user exists
	// Ensure test user exists
	testUser := model.User{
		Fullname: "Test User",
		Email:    "testuser@example.com",
		Password: "password123",  // Set password directly without hashing
	}
	result := db.Create(&testUser)
	if result.Error != nil {
		log.Fatalf("failed to create test user: %v", result.Error)
	} else {
		log.Printf("test user created: %v", testUser)
	}

	// Verify the test user exists
	var checkUser model.User
	if err := db.Where("email = ?", "testuser@example.com").First(&checkUser).Error; err != nil {
		log.Fatalf("failed to find test user: %v", err)
	} else {
		log.Printf("test user found: %v", checkUser)
	}
})

var _ = AfterSuite(func() {
	// Clean up the test database
	db.Exec("DROP SCHEMA public CASCADE")
	db.Exec("CREATE SCHEMA public")
})

var _ = Describe("UserAPI", func() {
	Describe("User Registration", func() {
    Context("when registering a new user", func() {
        It("should register a new user successfully", func() {
            // Prepare request body
            requestBody := map[string]interface{}{
                "fullname": "New User",
                "email":    "newuser@example.com",
                "password": "password123",
            }
            jsonBody, err := json.Marshal(requestBody)
            Expect(err).NotTo(HaveOccurred())

            // Perform HTTP request
            w := httptest.NewRecorder()
            req, err := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBuffer(jsonBody))
            Expect(err).NotTo(HaveOccurred())
            req.Header.Set("Content-Type", "application/json")
            router.ServeHTTP(w, req)

            // Check response status code
            Expect(w.Code).To(Equal(http.StatusCreated))

            // Check response body
            var response map[string]string
            err = json.Unmarshal(w.Body.Bytes(), &response)
            Expect(err).NotTo(HaveOccurred())
            Expect(response["message"]).To(Equal("register success"))
        })
    })
})

	Describe("User Login", func() {
		It("should login the user", func() {
			requestBody := map[string]interface{}{
				"email":    "newuser@example.com",
				"password": "password123",
			}
			jsonBody, err := json.Marshal(requestBody)
			Expect(err).NotTo(HaveOccurred())

			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonBody))
			Expect(err).NotTo(HaveOccurred())
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Set-Cookie")).To(ContainSubstring("session_token"))

			// Extract session token for further use
			cookies := w.Result().Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "session_token" {
					token = cookie.Value
				}
			}
		})
	})

	Describe("Get User Task Category", func() {
		BeforeEach(func() {
			// Ensure user is logged in and token is set
			if token == "" {
				requestBody := map[string]interface{}{
					"email":    "newuser@example.com",
					"password": "password123",
				}
				jsonBody, err := json.Marshal(requestBody)
				Expect(err).NotTo(HaveOccurred())

				w := httptest.NewRecorder()
				req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonBody))
				Expect(err).NotTo(HaveOccurred())
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "session_token" {
						token = cookie.Value
					}
				}
			}
		})

		It("should get tasks by category for the user", func() {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/v1/task/1", nil) // Assuming userID 1
			Expect(err).NotTo(HaveOccurred())
			req.Header.Set("Cookie", "session_token="+token)
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			var response []model.UserTaskCategory
			err = json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
			Expect(response).NotTo(BeNil())
		})
	})

	// Describe("User Logout", func() {
	// 	BeforeEach(func() {
	// 		// Ensure user is logged in and token is set
	// 		if token == "" {
	// 			requestBody := map[string]interface{}{
	// 				"email":    "testuser@example.com",
	// 				"password": "password123",
	// 			}
	// 			jsonBody, err := json.Marshal(requestBody)
	// 			Expect(err).NotTo(HaveOccurred())

	// 			w := httptest.NewRecorder()
	// 			req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonBody))
	// 			Expect(err).NotTo(HaveOccurred())
	// 			router.ServeHTTP(w, req)
	// 			Expect(w.Code).To(Equal(http.StatusOK))
	// 			cookies := w.Result().Cookies()
	// 			for _, cookie := range cookies {
	// 				if cookie.Name == "session_token" {
	// 					token = cookie.Value
	// 				}
	// 			}
	// 		}
	// 	})

	// 	It("should logout the user", func() {
	// 		w := httptest.NewRecorder()
	// 		req, err := http.NewRequest("POST", "/api/v1/user/logout", nil)
	// 		Expect(err).NotTo(HaveOccurred())
	// 		req.Header.Set("Cookie", "session_token="+token)
	// 		router.ServeHTTP(w, req)

	// 		Expect(w.Code).To(Equal(http.StatusOK))
	// 		var response map[string]string
	// 		err = json.Unmarshal(w.Body.Bytes(), &response)
	// 		Expect(err).NotTo(HaveOccurred())
	// 		Expect(response["message"]).To(Equal("logout success"))

	// 		// Verify session is removed from the database
	// 		session, err := sessionRepo.SessionAvailEmail("newuser@example.com")
	// 		Expect(err).To(HaveOccurred())
	// 		Expect(session.Email).To(BeEmpty())
	// 	})
	// })

	Describe("Task Management", func() {
    BeforeEach(func() {
        // Ensure user is logged in and token is set
        if token == "" {
            requestBody := map[string]interface{}{
                "email":    "newuser@example.com",
                "password": "password123",
            }
            jsonBody, err := json.Marshal(requestBody)
            Expect(err).NotTo(HaveOccurred())

            w := httptest.NewRecorder()
            req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonBody))
            Expect(err).NotTo(HaveOccurred())
            router.ServeHTTP(w, req)
            Expect(w.Code).To(Equal(http.StatusOK))
            cookies := w.Result().Cookies()
            for _, cookie := range cookies {
                if cookie.Name == "session_token" {
                    token = cookie.Value
                }
            }
        }
    })

    It("should create a new task", func() {
			w := httptest.NewRecorder()
			body := bytes.NewBufferString(`{"title": "New Task", "deadline": "2024-12-31T23:59:59Z", "priority": 1, "status": "pending", "category_id": 1, "user_id": 1}`)
			req, _ := http.NewRequest("POST", "/api/v1/tasks", body)
			req.Header.Set("Cookie", "session_token="+token)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
	
			Expect(w.Code).To(Equal(http.StatusCreated))
	
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())
	
			message, messageExists := response["message"]
			Expect(messageExists).To(BeTrue(), "Expected response to contain 'message' field")
			Expect(message).To(Equal("add task success"))
    })

    It("should update an existing task", func() {
        w := httptest.NewRecorder()
        body := bytes.NewBufferString(`{"title": "Updated Task", "description": "Updated description", "category_id": 1}`)
        req, _ := http.NewRequest("PUT", "/api/v1/tasks/1", body)  // Assuming taskID 1
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusOK))

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        Expect(err).NotTo(HaveOccurred())

        message, messageExists := response["message"]
        Expect(messageExists).To(BeTrue(), "Expected response to contain 'message' field")
        Expect(message).To(Equal("update task success"))
    })

    It("should delete an existing task", func() {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)  // Assuming taskID 1
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusOK))

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        Expect(err).NotTo(HaveOccurred())

        message, messageExists := response["message"]
        Expect(messageExists).To(BeTrue(), "Expected response to contain 'message' field")
        Expect(message).To(Equal("Task deleted successfully"))
    })

    It("should get a task by ID", func() {
			CreateTestTask(1, "Test Task", "2024-12-31T23:59:59Z", 1, "pending", 1, 1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/tasks/1", nil)  // Assuming taskID 1
			req.Header.Set("Cookie", "session_token="+token)
			router.ServeHTTP(w, req)
	
			Expect(w.Code).To(Equal(http.StatusOK))
	
			var task model.Task
			err := json.Unmarshal(w.Body.Bytes(), &task)
			Expect(err).NotTo(HaveOccurred())
			Expect(task.ID).To(Equal(1))
    })

    It("should get a list of tasks", func() {
			CreateTestTask(1, "Test Task", "2024-12-31T23:59:59Z", 1, "pending", 1, 1)

      w := httptest.NewRecorder()
      req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
      req.Header.Set("Cookie", "session_token="+token)
      router.ServeHTTP(w, req)

      Expect(w.Code).To(Equal(http.StatusOK))

      var tasks []model.Task
      err := json.Unmarshal(w.Body.Bytes(), &tasks)
      Expect(err).NotTo(HaveOccurred())
      Expect(tasks).NotTo(BeEmpty())
    })
	})

	Describe("Category Management", func() {
    BeforeEach(func() {
        // Ensure user is logged in and token is set
        if token == "" {
            requestBody := map[string]interface{}{
                "email":    "newuser@example.com",
                "password": "password123",
            }
            jsonBody, err := json.Marshal(requestBody)
            Expect(err).NotTo(HaveOccurred())

            w := httptest.NewRecorder()
            req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(jsonBody))
            Expect(err).NotTo(HaveOccurred())
            router.ServeHTTP(w, req)
            Expect(w.Code).To(Equal(http.StatusOK))
            cookies := w.Result().Cookies()
            for _, cookie := range cookies {
                if cookie.Name == "session_token" {
                    token = cookie.Value
                }
            }
        }
    })

    It("should create a new category", func() {
        w := httptest.NewRecorder()
        body := bytes.NewBufferString(`{"name": "Work"}`)
        req, _ := http.NewRequest("POST", "/api/v1/category", body)
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusCreated))

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        Expect(err).NotTo(HaveOccurred())

        message, messageExists := response["message"]
        Expect(messageExists).To(BeTrue(), "Expected response to contain 'message' field")
        Expect(message).To(Equal("add category success"))
    })

    It("should update an existing category", func() {
        w := httptest.NewRecorder()
        body := bytes.NewBufferString(`{"name": "Updated Category"}`)
        req, _ := http.NewRequest("PUT", "/api/v1/categories/1", body)  // Assuming categoryID 1
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusOK))

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        Expect(err).NotTo(HaveOccurred())

        message, messageExists := response["message"]
        Expect(messageExists).To(BeTrue(), "Expected response to contain 'message' field")
        Expect(message).To(Equal("category update success"))
    })

    It("should delete an existing category", func() {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("DELETE", "/api/v1/categories/1", nil)  // Assuming categoryID 1
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusOK))

        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        Expect(err).NotTo(HaveOccurred())

        message, messageExists := response["message"]
        Expect(messageExists).To(BeTrue(), "Expected response to contain 'message' field")
        Expect(message).To(Equal("Category deleted successfully"))
    })

    It("should get a category by ID", func() {
				CreateTestCategory(1, "Test Category")
			
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/categories/1", nil)  // Assuming categoryID 1
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusOK))

        var category model.Category
        err := json.Unmarshal(w.Body.Bytes(), &category)
        Expect(err).NotTo(HaveOccurred())
        Expect(category.ID).To(Equal(1))
    })

    It("should get a list of categories", func() {
			CreateTestCategory(1, "Test Category")
			
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/api/v1/categories", nil)
        req.Header.Set("Cookie", "session_token="+token)
        router.ServeHTTP(w, req)

        Expect(w.Code).To(Equal(http.StatusOK))

        var categories []model.Category
        err := json.Unmarshal(w.Body.Bytes(), &categories)
        Expect(err).NotTo(HaveOccurred())
        Expect(categories).NotTo(BeEmpty())
    })
	})
})
