package postgres

import (
    "errors"
    "fmt"

    "ecommerce-api/internal/domain/entity"
    domainrepo "ecommerce-api/internal/domain/repository"
    "ecommerce-api/internal/repository/models"
    "ecommerce-api/pkg/apperror"

    "gorm.io/gorm"
)

type orderRepo struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domainrepo.OrderRepository {
    return &orderRepo{db: db}
}

func (r *orderRepo) Create(o *entity.Order) error {
    m := models.EntityToOrderModel(o)
    if err := r.db.Create(m).Error; err != nil {
        return fmt.Errorf("orderRepo.Create: %w", err)
    }
    o.ID        = m.ID
    o.CreatedAt = m.CreatedAt
    o.UpdatedAt = m.UpdatedAt
    return nil
}

func (r *orderRepo) FindByID(id uint) (*entity.Order, error) {
    var m models.OrderModel
    err := r.db.Preload("Items").First(&m, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperror.ErrNotFound
        }
        return nil, fmt.Errorf("orderRepo.FindByID: %w", err)
    }
    return m.OrderModelToEntity(), nil
}

func (r *orderRepo) FindByUserID(userID uint, page, limit int) ([]entity.Order, int64, error) {
    var ms    []models.OrderModel
    var total int64
    offset := (page - 1) * limit

    r.db.Model(&models.OrderModel{}).Where("user_id = ?", userID).Count(&total)
    err := r.db.Preload("Items").
        Where("user_id = ?", userID).
        Offset(offset).Limit(limit).
        Find(&ms).Error
    if err != nil {
        return nil, 0, fmt.Errorf("orderRepo.FindByUserID: %w", err)
    }

    orders := make([]entity.Order, len(ms))
    for i, m := range ms {
        orders[i] = *m.OrderModelToEntity()
    }
    return orders, total, nil
}

func (r *orderRepo) UpdateStatus(id uint, status entity.OrderStatus) error {
    result := r.db.Model(&models.OrderModel{}).Where("id = ?", id).Update("status", string(status))
    if result.Error != nil {
        return fmt.Errorf("orderRepo.UpdateStatus: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return apperror.ErrNotFound
    }
    return nil
}