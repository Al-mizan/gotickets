package user

import (
	"errors"
	"fmt"
	"gotickets/internal/apperror"
	"gotickets/internal/auth"
	"gotickets/internal/domain/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

// Service defines the contract for user business logic.
// Implementations can be swapped or mocked in tests.
type Service interface {
	CreateUser(req dto.CreateRequest) (*dto.Response, error)
	LoginUser(req dto.LoginRequest) (*dto.Response, error)
	GetUserByID(userId uint) (*dto.Response, error)
}

type service struct {
	repo       Repository
	jwtService auth.JWTService
}

func NewService(repo Repository, jwtService auth.JWTService) Service {
	return &service{repo, jwtService}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {

	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	// hash password and set to user.Password
	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to hash password")
	}

	err = s.repo.CreateUser(&user)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return nil, apperror.NewConflict(err, "failed to create user: email already exists")
		}
		return nil, apperror.NewInternal(err, "failed to create user")
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil

}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user")
	}

	if user == nil {
		return nil, apperror.NewUnauthorized(ErrInvalidCredentials, "invalid email or password")
	}

	// check password
	err = user.checkPassword(req.Password)

	if err != nil {
		return nil, apperror.NewUnauthorized(ErrInvalidCredentials, "invalid email or password")
	}

	// generate token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to generate token")
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil
}

func (s *service) GetUserByID(userId uint) (*dto.Response, error) {
	user, err := s.repo.GetUserByID(userId)
	if err != nil {
		return nil, apperror.NewInternal(err, "failed to fetch user")
	}
	if user == nil {
		return nil, apperror.NewNotFound(fmt.Errorf("user not found"), "user not found")
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return &response, nil
}
