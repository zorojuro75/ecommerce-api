package postgres

import (
	"errors"
	"fmt"
	"strings"

	"ecommerce-api/internal/domain/entity"
	domainrepo "ecommerce-api/internal/domain/repository"
	"ecommerce-api/internal/repository/models"
	"ecommerce-api/pkg/apperror"

	"gorm.io/gorm"
)

type productRepo struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domainrepo.ProductRepository {
    return &productRepo{db: db}
}

func (r *productRepo) Create(p *entity.Product) error {
    m := models.ProductFromEntity(p)
    if err := r.db.Create(m).Error; err != nil {
        return fmt.Errorf("productRepo.Create: %w", err)
    }
    p.ID        = m.ID
    p.CreatedAt = m.CreatedAt
    p.UpdatedAt = m.UpdatedAt
    return nil
}

func (r *productRepo) FindByID(id uint) (*entity.Product, error) {
    var m models.ProductModel
    if err := r.db.First(&m, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperror.ErrNotFound
        }
        return nil, fmt.Errorf("productRepo.FindByID: %w", err)
    }
    return m.ToEntity(), nil
}

func (r *productRepo) FindAll(filter entity.ProductFilter) ([]entity.Product, int64, error) {
    var ms    []models.ProductModel
    var total int64

    query := r.db.Model(&models.ProductModel{})

    if filter.Search != "" {
        search := "%" + strings.ToLower(filter.Search) + "%"
        query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", search, search)
    }

    if filter.MinPrice > 0 {
        query = query.Where("price >= ?", filter.MinPrice)
    }
    if filter.MaxPrice > 0 {
        query = query.Where("price <= ?", filter.MaxPrice)
    }

    if err := query.Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("productRepo.FindAll count: %w", err)
    }

    switch filter.Sort {
    case "price_asc":  query = query.Order("price ASC")
    case "price_desc": query = query.Order("price DESC")
    case "name_asc":   query = query.Order("name ASC")
    case "newest":    query = query.Order("created_at DESC")
    default:          query = query.Order("id ASC")
    }
    offset := (filter.Page - 1) * filter.Limit
    if err := query.Offset(offset).Limit(filter.Limit).Find(&ms).Error; err != nil {
        return nil, 0, fmt.Errorf("productRepo.FindAll: %w", err)
    }

    products := make([]entity.Product, len(ms))
    for i, m := range ms {
        products[i] = *m.ToEntity()
    }
    return products, total, nil
}

func (r *productRepo) Update(p *entity.Product) error {
    m := models.ProductFromEntity(p)
    result := r.db.Save(m)
    if result.Error != nil {
        return fmt.Errorf("productRepo.Update: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return apperror.ErrNotFound
    }
    p.UpdatedAt = m.UpdatedAt
    return nil
}

func (r *productRepo) Delete(id uint) error {
    result := r.db.Delete(&models.ProductModel{}, id)
    if result.Error != nil {
        return fmt.Errorf("productRepo.Delete: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return apperror.ErrNotFound
    }
    return nil
}