package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"task-management/dao"
	"task-management/models"
	u "task-management/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Secret key for signing the JWT token (should be stored securely in an environment variable)
var jwtSecret []byte

func init() {
	if os.Getenv("token_password") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}
	jwtSecret = []byte(os.Getenv("token_password"))
}

func CreateAccount(c *gin.Context) {

	account := &models.Account{}
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp := account.Create() //Create account
	c.JSON(http.StatusCreated, resp)
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Ensure it follows "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Inside the callback function checks if the token uses HMAC signing.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signing method"})
				c.Abort()
			}
			// This is mainly for if statements.  When the app is running, this should already be initialized.
			if jwtSecret == nil {
				jwtSecret = []byte(os.Getenv("token_password"))
			} else if len(jwtSecret) == 0 {
				jwtSecret = []byte(os.Getenv("token_password"))
			}
			fmt.Println(jwtSecret)
			return jwtSecret, nil
		})
		if err != nil {
			println("Error: " + err.Error())
		}
		// If token is invalid
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid, continue request
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
		c.JSON(http.StatusUnauthorized, msg)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(loginRequest.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		msg = u.Message(false, "Invalid login credentials. Please try again")
		c.JSON(http.StatusUnauthorized, msg)
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	claims := jwt.MapClaims{
		"username": account.Email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	// Create a JWT token with an expiration time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		println("Token signing error: " + err.Error())
		c.Abort()
	}

	resp := u.Message(true, "Logged In")
	resp["token"] = tokenString
	c.JSON(http.StatusOK, resp)
}
