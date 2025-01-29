package controllers

import (
	"net/http"
	"task-management/dao"
	"task-management/models"

	"github.com/gin-gonic/gin"
)

func GetAllTasks(c *gin.Context) {
	var tasks []models.Task
	dao.GetDB().Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	var task models.Task
	if err := dao.GetDB().First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found!"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	// Bind and validate the request body to the task model
	if err := c.ShouldBindJSON(&task); err != nil {
		// If validation fails, return a 400 Bad Request with the error details
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional custom validation (e.g., checking status)
	if err := task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data", "details": err.Error()})
		return
	}

	// Save task to the database
	if err := dao.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task

	// Get task ID from URL parameters
	id := c.Param("id")

	// Find the task by ID
	if err := dao.GetDB().First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Bind the request body to the task model
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the task
	if err := task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data", "details": err.Error()})
		return
	}

	// Update the task in the database
	if err := dao.GetDB().Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Return the updated task
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	if err := dao.GetDB().Delete(&models.Task{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully!"})
}
