package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string
	Password    string
	Name        string
	Surname     string
	Birthdate   time.Time
	Active      bool
	LastLoginAt time.Time
}

type Query[T any] interface {
	// SELECT * FROM @@table WHERE id=@id
	GetByID(id int) (T, error)
}
