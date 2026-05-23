package user

import (
	"fmt"
	"gotickets/internal/user/dto"
)

var ErrInvalidCredentials = fmt.Errorf("invalid email or password")

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(req dto.CreateRequest) (*dto.Response, error) {

	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	// hash password and set to user.Password
	err := user.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// generate token

	err = s.repo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil

}

func (s *service) LoginUser(req dto.LoginRequest) (*dto.Response, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials // User not found
	}

	// check password
	err = user.checkPassword(req.Password)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	response := dto.Response{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}

	return &response, nil
}
