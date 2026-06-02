package main

//go:generate buf generate

import (
	"fmt"
	"log"
	"net/http"
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
	fmt.Printf("db.Name(): %v\n", db.Name())

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	// ROUTES
	routes.RouteUsers("/user", e, db)

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

	log.Println("Server is running on port :8080 with HTTP and gRPC support")
}
