package contract

import (
	"ecommerce-api/internal/domain/entity"
)

type UserUsecase interface {
    Register(req RegisterRequest) (*entity.User, error)
    Login(req LoginRequest) (string, error)
    GetUser(id uint) (*entity.User, error)
}

type RegisterRequest struct {
    Name     string
    Email    string
    Password string
}

type LoginRequest struct {
    Email    string
    Password string
}