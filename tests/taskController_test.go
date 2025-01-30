package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"task-management/dao"
	"task-management/models"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var testToken = generateTestToken()

// Test listing public tasks
func TestGetPublicTasks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/public/tasks", nil)
	resp := httptest.NewRecorder()
	GetTestRouter().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test creating a task (Protected Route)
func TestCreateTask(t *testing.T) {
	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
	}

	jsonValue, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/tasks/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+testToken) // Mock API key for authentication
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	GetTestRouter().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

// Test retrieving a task by ID
func TestGetTaskByID(t *testing.T) {
	// Create a task
	task := models.Task{Title: "Fetch Task", Description: "Fetch test", Status: "pending"}
	dao.GetDB().Create(&task)

	// Fetch the task
	req, _ := http.NewRequest("GET", "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)
	resp := httptest.NewRecorder()
	GetTestRouter().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test updating a task
func TestUpdateTask(t *testing.T) {
	// Create a task
	task := models.Task{Title: "Old Title", Description: "Old Desc", Status: "pending"}
	dao.GetDB().Create(&task)

	// Update data
	updateData := map[string]string{"title": "Updated Title"}
	jsonValue, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(int(task.ID)), bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	GetTestRouter().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test deleting a task
func TestDeleteTask(t *testing.T) {
	// Create a task
	task := models.Task{Title: "To be deleted", Description: "Delete me", Status: "pending"}
	dao.GetDB().Create(&task)

	// Delete request
	req, _ := http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	req.Header.Set("Authorization", "Bearer "+testToken)

	resp := httptest.NewRecorder()
	GetTestRouter().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test authentication middleware
func TestAuthMiddleware(t *testing.T) {
	req, _ := http.NewRequest("POST", "/tasks/", nil)
	req.Header.Set("Authorization", "wrong_api_key") // Invalid token

	resp := httptest.NewRecorder()
	GetTestRouter().ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func generateTestToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Expire in 24 hours
	})

	signedToken, _ := token.SignedString([]byte("your-secret-key")) // Use same key as middleware
	return signedToken
}
