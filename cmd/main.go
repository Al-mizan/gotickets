package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type User struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
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
	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/users", func(c *echo.Context) error {
		u := new(User)

		// binding the user data
		if err := c.Bind(u); err != nil {
			return err
		}

		// validating the user data
		if err := c.Validate(u); err != nil {
			return c.JSON()
		}

		// Process the user data (e.g., save to database)

		return c.JSON(http.StatusCreated, u)
		// or
		// return c.XML(http.StatusCreated, u)
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
