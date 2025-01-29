package controllers

import (
	"net/http"
	"os"
	"task-management/dao"
	"task-management/models"
	u "task-management/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateAccount(c *gin.Context) {

	account := &models.Account{}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp := account.Create() //Create account
	c.JSON(http.StatusCreated, resp)
}

// Secret key for signing the JWT token (should be stored securely in an environment variable)
var jwtSecret = []byte("your-secret-key")

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer your-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	account := &models.Account{}
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var msg map[string]interface{}
	err := dao.GetDB().Table("accounts").Where("email = ?", loginRequest.Email).First(account).Error
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			msg = u.Message(false, "Email address not found")
		} else {
			msg = u.Message(false, "Connection error. Please retry")
		}
		c.JSON(http.StatusForbidden, msg)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginRequest.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		msg = u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &models.Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	c.JSON(http.StatusOK, resp)
}
