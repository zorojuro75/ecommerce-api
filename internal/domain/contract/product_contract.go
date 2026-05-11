package contract

import (
	"ecommerce-api/internal/domain/entity"
)

type ProductUsecase interface {
    CreateProduct(req CreateProductRequest) (*entity.Product, error)
    GetProduct(id uint) (*entity.Product, error)
    ListProducts(page, limit int) ([]entity.Product, int64, error)
    UpdateProduct(id uint, req UpdateProductRequest) (*entity.Product, error)
    DeleteProduct(id uint) error
}

type CreateProductRequest struct {
    Name        string
    Description string
    Price       float64
    Stock       int
    CategoryID  uint
}

type UpdateProductRequest struct {
    Name        string
    Description string
    Price       float64
    Stock       int
    CategoryID  uint
}