package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-management/controllers"
	"task-management/dao"
	"task-management/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	dao.SetDb(db)
	dao.GetDB().AutoMigrate(&models.Account{})
	router := gin.Default()
	router.POST("/login", controllers.Login)
	return router
}

func TestLoginSuccess(t *testing.T) {
	router := setupAuthRouter()

	// Create an account
	acct := models.Account{Email: "admin@taskmgmt.com", Password: "login test"}
	dao.GetDB().Create(&acct)

	loginData := controllers.LoginRequest{
		Email:    "admin@taskmgmt.com",
		Password: "login test",
	}

	loginJSON, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var responseData map[string]string
	json.Unmarshal(resp.Body.Bytes(), &responseData)

	assert.NotEmpty(t, responseData["token"])
}

func TestLoginFailure(t *testing.T) {
	router := setupAuthRouter()

	loginData := controllers.LoginRequest{
		Email:    "wronguser@gmail.com",
		Password: "password123",
	}

	loginJSON, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
