package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "user-ms/docs" // Import generated docs
	"user-ms/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Setup database connection and logging
func setupDatabase() {

	err2 := godotenv.Load()
	if err2 != nil {
		log.Fatalf("Error loading .env file: %v", err2)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("Error: Missing database environment variables in .env file.")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host,
		user,
		password,
		dbname,
		port,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Info, // Log level
			Colorful:                  true,        // Disable color
			IgnoreRecordNotFoundError: true,
		},
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.UserEnity{}) //Auto migrate the product model
}

// @Summary Create a user
// @Description Creates a new user.
// @Accept json
// @Produce json
// @Param request body models.UserCreationRequest true "User creation request"
// @Success 200 {object} models.UserCreationResposnce "User creation response"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /users [post]
func createUser(c *gin.Context) {
	var request models.UserCreationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	// Adapt from request to model (assuming you have a separate User model)
	user := models.UserEnity{
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password, // In real app, hash this!
		Gender:   request.Gender,
	}

	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to create user",
			Error:   result.Error.Error(),
		})
		return
	}

	// Construct the response
	response := models.UserCreationResposnce{
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password, //  Don't return plain password in real app
		Gender:   user.Gender,
		ID:       int16(user.ID), // Convert uint to int16
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get all users
// @Description Retrieves all users
// @Produce json
// @Success 200 {array} models.UserListResponse
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /users [get]
func getUsers(c *gin.Context) {
	var users []models.UserEnity // Use models.UserEntity
	result := db.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Failed to get users",
			Error:   result.Error.Error(),
		})
		return
	}

	// Parse the data into the UserListResponse struct
	var response []models.UserListResponse
	for _, user := range users {
		response = append(response, models.UserListResponse{
			UserName: user.UserName,
			Email:    user.Email,
			Gender:   user.Gender,
			ID:       int16(user.ID), // Convert uint to int16
		})
	}

	c.JSON(http.StatusOK, response)
}

// @title Product API
// @version 1.0
// @description This is a sample service for managing products.
// @host localhost:8080
// @BasePath /
func main() {
	setupDatabase()
	r := gin.Default()

	r.POST("/users", createUser)
	r.GET("/users", getUsers)

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	logrus.Info("Starting server on port 8080")
	r.Run(":8080")
}
