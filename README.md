# Go Rest Task Manager

## Project Description

Go Task Manager is a comprehensive task management application developed using the Go programming language. This application empowers users to efficiently organize their tasks by providing features for task and category management. Additionally, it includes user management utilities for user registration, login, logout, and fetching tasks by user id.

## Tech Stack

- **Programming Language**: Go (Golang)
- **Database**: PostgreSQL
- **ORM**: Gorm
- **API Framework**: Gin
- **Containerization**: Docker
- **Caching**: Redis
- **Authentication**: JSON Web Token (JWT)
- **Testing**: Ginkgo

## Features

- **User Management**: Utility features for user management, including registration, login, and logout.
- **Task Management**: CRUD operations for tasks, including creation, retrieval, updating, and deletion.
- **Category Management**: CRUD operations for categories to organize tasks effectively.
- **Authentication**: JWT-based authentication for API security.
- **Input Validation**: Validation of user input data.
- **API Documentation**: Includes a Postman Collection for easy testing and interaction with the API endpoints.
- **Testing**: Comprehensive unit testing using Ginkgo to ensure the reliability and correctness of the API endpoints.
- **Caching**: Utilizes Redis to cache task list and category list for improved performance.

## How to Run the Program

### Prerequisites

Ensure you have installed:

- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/dl/)
- [Ginkgo](https://github.com/onsi/ginkgo) (Ginkgo is required for running unit tests.)

### Steps

1. **Pull Docker Image**

```bash
docker pull daffajatmiko/go-task-manager:v1.0.1
```

1. **Run Docker Containers**

```bash
docker run -p 8080:8080 daffajatmiko/go-task-manager:v1.0.1
```

3. **Accessing the API**

The application will run at http://localhost:8080. You can access the application through your web browser or postman.

4. **Manaing the Database**

To manage the PostgreSQL database used in this project, you can utilize pgAdmin, a web-based PostgreSQL administration tool. Follow these steps to access and manage the database:

1. Open a web browser and go to http://localhost:8082.
2. Log in with the following credentials:
   - Email: admin@admin.com
   - Password: admin
3. Once logged in, you can register a PostgreSQL server with details matching your Docker container configuration.
4. Access the database and its tables to view, insert, update, or delete data as needed.

## Testing

This project includes comprehensive unit testing to ensure the reliability and correctness of the API endpoints. To run the tests using Ginkgo, follow these steps:

1. Navigate to the service folder in your terminal.
2. Run the following command to execute the tests:

```bash
go test
/ or ginkgo

```
