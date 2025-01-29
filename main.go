package main

import (
	"log"
	"task-management/dao"
	"task-management/models"
	"task-management/routes" // Assuming your SetupRoutes function is in this package

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	dao.ConnectDatabase()

	// Apply migrations to create/update the database schema
	dao.GetDB().AutoMigrate(&models.Task{})
	dao.GetDB().AutoMigrate(&models.Account{})
	dao.GetDB().AutoMigrate(&models.Token{})

	// Initialize the Gin router
	router := gin.Default()

	// Set up routes for the API
	routes.SetupRoutes(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
