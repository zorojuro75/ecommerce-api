package entity

import (
    "strings"
    "time"
    "ecommerce-api/pkg/apperror"
)

type Product struct {
    ID          uint
    Name        string
    Description string
    Price       float64
    Stock       int
    CategoryID  uint
    CreatedAt   time.Time
    UpdatedAt   time.Time
}


func (p Product) Validate() error {
    if strings.TrimSpace(p.Name) == "" {
        return apperror.ErrInvalidInput
    }
    if p.Price <= 0 {
        return apperror.ErrInvalidInput
    }
    if p.Stock < 0 {
        return apperror.ErrInvalidInput
    }
    return nil
}

func (p Product) IsAvailable() bool {
    return p.Stock > 0
}

type ProductUsecase interface {
    CreateProduct(req CreateProductRequest) (*Product, error)
    GetProduct(id uint) (*Product, error)
    ListProducts(page, limit int) ([]Product, int64, error)
    UpdateProduct(id uint, req UpdateProductRequest) (*Product, error)
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