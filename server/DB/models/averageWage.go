package models

import (
	"gorm.io/gorm"
)

type AverageWage struct {
	gorm.Model
	Voivodeship string
	Year        int
	AverageWage float64
}
