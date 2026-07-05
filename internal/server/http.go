package server

import (
	"errors"
	"fmt"
	"gotickets/internal/apperror"
	"gotickets/internal/auth"
	"gotickets/internal/config"
	"gotickets/internal/domain/booking"
	"gotickets/internal/domain/event"
	"gotickets/internal/domain/user"
	"gotickets/internal/httpresponse"
	"gotickets/internal/middlewares"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

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

func Start(db *gorm.DB, cfg *config.Config) {
	if cfg.Environment == "development" {
		db.AutoMigrate(&user.User{}, &event.Event{}, &booking.Booking{})
	}

	e := echo.New()

	e.HTTPErrorHandler = func(c *echo.Context, err error) {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, httpresponse.Error{
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			return
		}

		var he *echo.HTTPError
		if errors.As(err, &he) {
			c.JSON(he.Code, httpresponse.Error{
				Code:    he.Code,
				Message: fmt.Sprintf("%v", he.Message),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
	}

	// global validator which is validate http request body
	e.Validator = &CustomValidator{validator: validator.New()}

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogURI:       true,
			LogMethod:    true,
			LogStatus:    true,
			LogLatency:   true,
			LogRemoteIP:  true,
			LogUserAgent: true,

			LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
				logger.Info("request",
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					"remote_ip", v.RemoteIP,
					"user_agent", v.UserAgent,
				)
				return nil
			},
		},
	))
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
	}))

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "running")
	})

	// Setup dependencies
	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		slog.Error("failed to initialize jwt service", "error", err)
		os.Exit(1)
	}
	authMw := middlewares.AuthMiddleware(jwtService)

	// User domain
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, jwtService)
	userHandler := user.NewHandler(userSvc)
	user.RegisterRoutes(e, userHandler, authMw)

	// Event domain
	eventRepo := event.NewRepository(db)
	eventSvc := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventSvc)
	event.RegisterRoutes(e, eventHandler, authMw)

	// Booking domain
	bookingRepo := booking.NewRepository(db)
	bookingSvc := booking.NewService(bookingRepo, eventRepo)
	bookingHandler := booking.NewHandler(bookingSvc)
	booking.RegisterRoutes(e, bookingHandler, authMw)

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
