package repository

import "ecommerce-api/internal/domain/entity"

type OrderRepository interface {
    Create(o *entity.Order) error
    FindByID(id uint) (*entity.Order, error)
    FindByUserID(userID uint, page, limit int) ([]entity.Order, int64, error)
    UpdateStatus(id uint, status entity.OrderStatus) error
}