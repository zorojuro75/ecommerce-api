package entity

import (
    "strings"
    "time"
    "ecommerce-api/pkg/apperror"
)

type UserRole string

const (
    RoleCustomer UserRole = "customer"
    RoleAdmin    UserRole = "admin"
)

type User struct {
    ID           uint
    Name         string
    Email        string
	Mobile		 string
    PasswordHash string
    Role         UserRole
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (u User) Validate() error {
    if strings.TrimSpace(u.Name) == "" {
        return apperror.ErrInvalidInput
    }
    if !strings.Contains(u.Email, "@") {
        return apperror.ErrInvalidInput
    }
    return nil
}

func (u User) IsAdmin() bool {
    return u.Role == RoleAdmin
}

type UserUsecase interface {
    Register(req RegisterRequest) (*User, error)
    Login(req LoginRequest) (string, error)
    GetUser(id uint) (*User, error)
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