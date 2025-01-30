package tests_test

import (
	"task-management/controllers"
	"task-management/dao"
	"task-management/models"
	"task-management/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testRouter *gin.Engine

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
	dao.GetDB().AutoMigrate(&models.Token{})
	// Setup routes
	routes.SetupRoutes(testRouter)

	protected := testRouter.Group("/")
	protected.Use(controllers.Authenticate())

	// Run tests
	m.Run()
}

// Expose a function to get the test router for unit tests
func GetTestRouter() *gin.Engine {
	return testRouter
}
