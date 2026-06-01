package routes

import (
	"net/http"
	"study4cash/DB/models"
	"study4cash/auth"
	"time"

	"github.com/alexedwards/argon2id"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RouteUsers(prefix string, e *echo.Echo, db *gorm.DB) {
	userRouter := e.Group(prefix)
	userRouter.POST("/register", func(c *echo.Context) error { return Register(c, db) })
	userRouter.POST("/login", func(c *echo.Context) error { return Login(c, db) })

	protected := userRouter.Group("")
	protected.Use(auth.JWTMiddleware)
	protected.GET("/details", func(c *echo.Context) error { return Details(c, db) })
}

type RegisterRequest struct {
	Email     string    `json:"email"  validate:"required,email"`
	Password  string    `json:"password"  validate:"required,min=8,max=32"`
	Name      string    `json:"name"  validate:"required,min=3,max=60"`
	Surname   string    `json:"surname"  validate:"required,min=3,max=60"`
	Birthdate time.Time `json:"birthdate"  validate:"required"`
}
type RegisterResponse struct {
	JWTToken string `json:"token"`
	UserID   uint   `json:"user_id"`
}

func Register(c *echo.Context, db *gorm.DB) error {
	registerData := new(RegisterRequest)
	if err := c.Bind(registerData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(registerData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := gorm.G[models.User](db).Where("email = ?", registerData.Email).Find(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(users) > 0 {
		return c.JSON(http.StatusBadRequest, "User already exists")
	}

	password, err := argon2id.CreateHash(registerData.Password, argon2id.DefaultParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	userModel := models.User{
		Email:       registerData.Email,
		Password:    password,
		Name:        registerData.Name,
		Surname:     registerData.Surname,
		Birthdate:   registerData.Birthdate,
		Active:      true,
		LastLoginAt: time.Now(),
	}
	tx := db.Save(&userModel)
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, tx.Error.Error())
	}

	token, err := auth.GenerateJWT(userModel.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, RegisterResponse{
		JWTToken: *token,
		UserID:   userModel.ID,
	})
}

type LoginRequest struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8,max=32"`
}
type LoginResponse struct {
	JWTToken string `json:"token"`
	UserID   uint   `json:"user_id"`
}

func Login(c *echo.Context, db *gorm.DB) error {
	loginData := new(RegisterRequest)
	if err := c.Bind(loginData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(loginData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	users, err := gorm.G[models.User](db).Where("email = ?", loginData.Email).Find(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(users) == 0 {
		return c.JSON(http.StatusBadRequest, "User not found")
	}

	user := users[0]
	match, err := argon2id.ComparePasswordAndHash(loginData.Password, user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not verify password")
	}
	if !match {
		return c.JSON(http.StatusUnauthorized, "Invalid password")
	}

	token, err := auth.GenerateJWT(user.ID)

	return c.JSON(http.StatusOK, LoginResponse{
		JWTToken: *token,
		UserID:   user.ID,
	})
}

type DetailsResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"  validate:"required,email"`
	Name      string    `json:"name"  validate:"required,min=3,max=60"`
	Surname   string    `json:"surname"  validate:"required,min=3,max=60"`
	Birthdate time.Time `json:"birthdate"  validate:"required"`
}

func Details(c *echo.Context, db *gorm.DB) error {
	userID := c.Get("userID").(uint)

	users, err := gorm.G[models.User](db).Where("id = ?", userID).Find(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(users) != 1 {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user := users[0]
	return c.JSON(http.StatusOK, DetailsResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		Birthdate: user.Birthdate,
	})
}
