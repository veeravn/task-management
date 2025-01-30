package routes

import (
	"task-management/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)
	router.POST("/account", controllers.CreateAccount)
	// Public routes (no authentication required)
	public := router.Group("/public")
	{
		public.GET("/tasks", controllers.GetAllTasks) // List all tasks
	}

	// Protected routes (authentication required)
	protected := router.Group("/tasks")
	protected.Use(controllers.Authenticate()) // Apply authentication middleware
	{
		protected.GET("/:id", controllers.GetTaskByID)   // Get task by ID
		protected.POST("/", controllers.CreateTask)      // Create new task
		protected.PUT("/:id", controllers.UpdateTask)    // Update task
		protected.DELETE("/:id", controllers.DeleteTask) // Delete task
	}
}
