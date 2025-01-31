package models

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

// Define Golang ENUM for Task Status
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in-progress"
	StatusCompleted  TaskStatus = "completed"
)

// Validate if TaskStatus is valid
func (s TaskStatus) IsValid() error {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return nil
	}
	return errors.New("invalid status")
}

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"unique;not null" json:"title" binding:"required,min=3,max=255"`
	Description string     `gorm:"not null" json:"description" binding:"required,min=5,max=500"`
	Status      TaskStatus `gorm:"type:text;default:'pending'" json:"status"` // Use TEXT instead of ENUM
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// This method can be used to perform custom validations on the task model before saving it to the database.
func (task *Task) Validate() error {
	validate := validator.New()

	// Example: You could validate fields like the status here, or even other custom business logic
	if err := validate.Struct(task); err != nil {
		return err
	}

	return nil
}
