package models

import (
	"gorm.io/gorm"
)

type Graduating struct {
	gorm.Model
	Voivodeship        string
	Year               int
	StudentsGraduating int
}
