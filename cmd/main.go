package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" validate:"required,email" gorm:"type:varchar(255);uniqueIndex;not nll"`
	Password string `json:"password" validate:"required,min=6" gorm:"type:varchar(100);not null"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func main() {

	dsn := "host=localhost user=postgres password=postgres dbname=gotickets port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		panic("failed to connect database")
	} else {
		println("Database connection successful")
	}

	db.AutoMigrate(&User{})

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Validator = &CustomValidator{validator: validator.New()}

	e.POST("/users", func(c *echo.Context) error {
		newUser := new(User)

		// binding the user data
		if err := c.Bind(newUser); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		}

		// validating the user data
		if err := c.Validate(newUser); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		}

		// save to database
		result := db.Create(newUser)
		if result.Error != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": result.Error.Error()})
		}

		return c.JSON(http.StatusCreated, newUser)
		// or
		// return c.XML(http.StatusCreated, u)
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
