package models

import (
	"strings"
	"task-management/dao"
	u "task-management/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

// a struct to rep task
type Task struct {
	gorm.Model
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func init() {
	dao.GetDB().Debug().AutoMigrate(&Task{})
}

// Validate incoming user details...
func (task *Task) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(task.Title, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if !strings.Contains(task.Description, "@") {
		return u.Message(false, "Description is required"), false
	}

	//Email must be unique
	temp := &Task{}

	//check for errors and duplicate tasks
	err := dao.GetDB().Table("tasks").Where("id = ?", task.Id).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Id != "" {
		return u.Message(false, "Id already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}
