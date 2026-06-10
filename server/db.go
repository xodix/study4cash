package main

import (
	"fmt"
	"log"
	"os"

	"study4cash/DB/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func configDB() *gorm.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	log.Println("DB_USERNAME:", username)
	log.Println("DB_PASSWORD:", password)
	log.Println("DB_HOST:", host)
	log.Println("DB_NAME:", DBName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, DBName)
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err.Error())
	}
	db.AutoMigrate(&models.User{})

	return db
}
