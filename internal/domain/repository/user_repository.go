package repository

import "ecommerce-api/internal/domain/entity"

type UserRepository interface {
    Create(u *entity.User) error
    FindByID(id uint) (*entity.User, error)
    FindByEmail(email string) (*entity.User, error)
    Update(u *entity.User) error
}