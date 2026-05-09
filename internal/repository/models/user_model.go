package models

import (
    "ecommerce-api/internal/domain/entity"
    "gorm.io/gorm"
)

type UserModel struct {
    gorm.Model
    Name         string            `gorm:"not null"`
    Email        string            `gorm:"uniqueIndex;not null"`
    PasswordHash string            `gorm:"not null"`
    Role         string            `gorm:"default:'customer'"`
}

func (UserModel) TableName() string { return "users" }

func (m *UserModel) ToEntity() *entity.User {
    return &entity.User{
        ID: m.ID, Name: m.Name, Email: m.Email,
        PasswordHash: m.PasswordHash,
        Role:         entity.UserRole(m.Role),
        CreatedAt:    m.CreatedAt, UpdatedAt: m.UpdatedAt,
    }
}

func UserFromEntity(u *entity.User) *UserModel {
    return &UserModel{
        Model:        gorm.Model{ID: u.ID},
        Name:         u.Name, Email: u.Email,
        PasswordHash: u.PasswordHash,
        Role:         string(u.Role),
    }
}