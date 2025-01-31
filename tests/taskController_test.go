package tests_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"task-management/controllers"
	"task-management/dao"
	"task-management/models"
	"task-management/routes"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testToken = generateTestToken()
var testRouter *gin.Engine

const tokenKey = "token_password"
const tokenVal = "thisIsTheJwtPassword"

// Setup the test database and router
func TestMain(m *testing.M) {
	// Initialize Gin router
	gin.SetMode(gin.TestMode)
	testRouter = gin.Default()

	// Setup an in-memory SQLite DB
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	dao.SetDb(db)
	dao.GetDB().AutoMigrate(&models.Task{})
	dao.GetDB().AutoMigrate(&models.Account{})
	// Setup routes
	routes.SetupRoutes(testRouter)

	protected := testRouter.Group("/")
	protected.Use(controllers.Authenticate())

	// Run tests
	m.Run()
}

// Helper function to create a valid test task
func createTestTask() models.Task {
	task := models.Task{
		Title:       "Valid Task",
		Description: "A valid task description",
		Status:      "pending",
	}
	dao.GetDB().Create(&task)
	return task
}

// Test listing public tasks
func TestGetPublicTasks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/public/tasks", nil)
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test creating a task (Protected Route)
func TestCreateTask(t *testing.T) {
	os.Setenv(tokenKey, tokenVal)
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
	testRouter.ServeHTTP(resp, req)

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
	testRouter.ServeHTTP(resp, req)

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
	testRouter.ServeHTTP(resp, req)

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
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test authentication middleware
func TestAuthMiddleware(t *testing.T) {
	req, _ := http.NewRequest("POST", "/tasks/", nil)
	req.Header.Set("Authorization", "wrong_api_key") // Invalid token

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// ❌ **Test: Create Task with Missing Fields**
func TestCreateTask_MissingFields(t *testing.T) {
	task := models.Task{
		Title: "", // Title is required
	}
	taskJSON, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/tasks/", bytes.NewBuffer(taskJSON))
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ❌ **Test: Create Duplicate Task**
func TestCreateTask_Duplicate(t *testing.T) {
	task := createTestTask() // Insert a task into DB

	// Try to create a task with the same title
	duplicateTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      "pending",
	}
	taskJSON, _ := json.Marshal(duplicateTask)

	req, _ := http.NewRequest("POST", "/tasks/", bytes.NewBuffer(taskJSON))
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code) // Expecting 409 Conflict
}

// ❌ **Test: Create Task Without Authorization**
func TestCreateTask_Unauthorized(t *testing.T) {
	task := models.Task{
		Title:       "Unauthorized Task",
		Description: "Testing unauthorized access",
		Status:      "pending",
	}
	taskJSON, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/tasks/", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json") // No Auth Header

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// ❌ **Test: Update Non-Existing Task**
func TestUpdateTask_NotFound(t *testing.T) {
	taskUpdate := models.Task{
		Title:       "Non-Existent Task",
		Description: "Trying to update a non-existent task",
		Status:      "completed",
	}
	taskJSON, _ := json.Marshal(taskUpdate)

	req, _ := http.NewRequest("PUT", "/tasks/99999", bytes.NewBuffer(taskJSON)) // Non-existent ID
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ❌ **Test: Update Task with Invalid Data**
func TestUpdateTask_InvalidData(t *testing.T) {

	taskUpdate := models.Task{
		Title:       "", // Empty title (invalid)
		Description: "Updated description",
		Status:      "completed",
	}
	taskJSON, _ := json.Marshal(taskUpdate)

	req, _ := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(1), bytes.NewBuffer(taskJSON))
	req.Header.Set("Authorization", "Bearer "+testToken)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// ❌ **Test: Update Task Without Authorization**
func TestUpdateTask_Unauthorized(t *testing.T) {
	task := createTestTask() // Insert a test task

	taskUpdate := models.Task{
		Title:       "Unauthorized Update",
		Description: "Testing unauthorized update",
		Status:      "completed",
	}
	taskJSON, _ := json.Marshal(taskUpdate)

	req, _ := http.NewRequest("PUT", "/tasks/"+strconv.Itoa(int(task.ID)), bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json") // No Auth Header

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// ❌ **Test: Delete Non-Existent Task**
func TestDeleteTask_NotFound(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/tasks/99999", nil) // Non-existent ID
	req.Header.Set("Authorization", "Bearer "+testToken)

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// ❌ **Test: Delete Task Without Authorization**
func TestDeleteTask_Unauthorized(t *testing.T) {
	task := createTestTask() // Insert a test task

	req, _ := http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(int(task.ID)), nil) // No Auth Header

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func generateTestToken() string {
	os.Setenv(tokenKey, tokenVal)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Expire in 24 hours
	})

	signedToken, _ := token.SignedString([]byte("thisIsTheJwtPassword")) // Use same key as middleware
	return signedToken
}
