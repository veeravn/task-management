package models

import "github.com/go-playground/validator/v10"

// a struct to rep task
type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" binding:"required,min=3,max=255"`                        // Title should be between 3 and 255 characters
	Description string `json:"description" binding:"required,min=5,max=500"`                  // Description should be between 5 and 500 characters
	Status      string `json:"status" binding:"required,oneof=pending in-progress completed"` // Status must be one of these values
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
