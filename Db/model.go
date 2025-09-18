package Db

import "time"

// Users table
type Users struct {
	UserID       int       `gorm:"primaryKey;column:user_id"`
	UserName     string    `gorm:"column:user_name"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	Role         string    `gorm:"column:role"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

// Category table
type Category struct {
	CategoryID   int    `gorm:"primaryKey;column:category_id"`
	CategoryName string `gorm:"column:category_name"`
	Description  string `gorm:"column:description"`
}

// Task table
type Task struct {
	TaskID      int       `gorm:"primaryKey;column:task_id"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	DueDate     time.Time `gorm:"column:due_date"`
	Status      string    `gorm:"column:status"`
	UserID      int       `gorm:"column:user_id"`
	CategoryID  int       `gorm:"column:category_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`

	Category Category `gorm:"foreignKey:category_id;references:category_id"`
}

// Log table
type Log struct {
	LogID     int       `gorm:"primaryKey;column:log_id"`
	UserID    int       `gorm:"column:user_id"`
	Action    string    `gorm:"column:action"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
