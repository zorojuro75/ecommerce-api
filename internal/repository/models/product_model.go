package models

import (
    "ecommerce-api/internal/domain/entity"
    "gorm.io/gorm"
)

type ProductModel struct {
    gorm.Model
    Name        string  `gorm:"not null;index"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"not null"`
    Stock       int     `gorm:"default:0"`
    CategoryID  uint    `gorm:"not null;index"`
}

func (ProductModel) TableName() string { return "products" }

// ToEntity — GORM model → domain entity (called when reading from DB)
func (m *ProductModel) ToEntity() *entity.Product {
    return &entity.Product{
        ID:          m.ID,
        Name:        m.Name,
        Description: m.Description,
        Price:       m.Price,
        Stock:       m.Stock,
        CategoryID:  m.CategoryID,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
    }
}

// FromEntity — domain entity → GORM model (called when writing to DB)
func ProductFromEntity(p *entity.Product) *ProductModel {
    return &ProductModel{
        Model:       gorm.Model{ID: p.ID},
        Name:        p.Name,
        Description: p.Description,
        Price:       p.Price,
        Stock:       p.Stock,
        CategoryID:  p.CategoryID,
    }
}