package models

import (
	"gorm.io/gorm"
)

type Attending struct {
	gorm.Model
	Voivodeship       string
	Year              int
	StudentsAttending int
}
