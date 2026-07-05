package user

import (
	"gotickets/internal/apperror"
	"gotickets/internal/ctxkeys"
	"gotickets/internal/domain/user/dto"

	"github.com/labstack/echo/v5"
)

// Handler defines the contract for user HTTP handlers.
// Implementations can be swapped or mocked in tests.
type Handler interface {
	CreateUser(c *echo.Context) error
	LoginUser(c *echo.Context) error
	GetMe(c *echo.Context) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest // input

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.CreateUser(req)
	if err != nil {
		return err
	}

	return c.JSON(201, response)

}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest // input

	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest(err, "Invalid request payload")
	}

	if err := c.Validate(&req); err != nil {
		return apperror.NewBadRequest(err, "Validation failed")
	}

	response, err := h.service.LoginUser(req)

	if err != nil {
		return err
	}

	return c.JSON(200, response)

}

func (h *handler) GetMe(c *echo.Context) error {
	userID, ok := c.Get(string(ctxkeys.UserID)).(uint)
	if !ok {
		return apperror.NewUnauthorized(nil, "missing user id in context")
	}

	response, err := h.service.GetUserByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(200, response)
}
