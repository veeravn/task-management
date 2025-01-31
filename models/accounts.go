package models

import (
	"errors"
	"fmt"
	"strings"
	"task-management/dao"
	u "task-management/utils"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := dao.GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if err != nil {
		rawError := fmt.Sprintf("%s", err)
		if rawError != gorm.ErrRecordNotFound.Error() {
			// Handle unexpected errors
			fmt.Println("An unexpected error occurred, raw:", errors.Unwrap(gorm.ErrRecordNotFound).Error())
			return u.Message(false, "An unexpected error occurred:"), false
		} else {
			println(rawError)
			return u.Message(false, "Requirement passed"), true
		}
	} else {
		// Continue with execution if a record is found
		return u.Message(false, "Requirement passed"), true
	}
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	dao.GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func GetUser(u uint) *Account {

	acc := &Account{}
	dao.GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
