package main

//go:generate buf generate

import (
	"fmt"
	"log"
	"net/http"
	"study4cash/DB/models"
	"study4cash/routes"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	godotenv.Load(".env")
	db := configDB()
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalln(err.Error())
	}
	if err := db.AutoMigrate(&models.Attending{}); err != nil {
		log.Fatalln(err.Error())
	}
	if err := db.AutoMigrate(&models.AverageWage{}); err != nil {
		log.Fatalln(err.Error())
	}
	if err := db.AutoMigrate(&models.Graduating{}); err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("db.Name(): %v\n", db.Name())

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// ROUTES
	routes.RouteUsers("/user", e, db)
	routes.RouteData("/data", e, db)

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

	log.Println("Server is running on port :8080 with HTTP and gRPC support")
}
