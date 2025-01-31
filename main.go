package main

import (
	"fmt"
	"log"
	"task-management/dao"
	"task-management/models"
	"task-management/routes" // Assuming your SetupRoutes function is in this package

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	// Connect to the database
	dao.ConnectDatabase()

	// Apply migrations to create/update the database schema
	// Run migrations
	dao.GetDB().AutoMigrate(&models.Task{})
	dao.GetDB().AutoMigrate(&models.Account{})

	// Initialize the Gin router
	router := gin.Default()

	// Set up routes for the API
	routes.SetupRoutes(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
