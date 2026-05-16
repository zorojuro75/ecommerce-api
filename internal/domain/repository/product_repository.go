package repository

import "ecommerce-api/internal/domain/entity"

type ProductRepository interface {
    Create(p *entity.Product) error
    FindByID(id uint) (*entity.Product, error)
    FindAll(filter entity.ProductFilter) ([]entity.Product, int64, error)
    Update(p *entity.Product) error
    Delete(id uint) error
}