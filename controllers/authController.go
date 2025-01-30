package controllers

import (
	"net/http"
	"os"
	"strings"
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
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// error := gin.Error{Err: bcrypt.ErrMismatchedHashAndPassword, Type: gin.ErrorTypePublic, Meta: "Invalid signing method"}
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signing method"})
				c.Abort()
			}
			return jwtSecret, nil
		})

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
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := &models.Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["token"] = account.Token
	c.JSON(http.StatusOK, resp)
}
