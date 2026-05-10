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

type userRepo struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainrepo.UserRepository {
    return &userRepo{db: db}
}

func (r *userRepo) Create(u *entity.User) error {
    m := models.UserFromEntity(u)
    if err := r.db.Create(m).Error; err != nil {
        return fmt.Errorf("userRepo.Create: %w", err)
    }
    u.ID        = m.ID
    u.CreatedAt = m.CreatedAt
    u.UpdatedAt = m.UpdatedAt
    return nil
}

func (r *userRepo) FindByID(id uint) (*entity.User, error) {
    var m models.UserModel
    if err := r.db.First(&m, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperror.ErrNotFound
        }
        return nil, fmt.Errorf("userRepo.FindByID: %w", err)
    }
    return m.ToEntity(), nil
}

func (r *userRepo) FindByEmail(email string) (*entity.User, error) {
    var m models.UserModel
    if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperror.ErrNotFound
        }
        return nil, fmt.Errorf("userRepo.FindByEmail: %w", err)
    }
    return m.ToEntity(), nil
}

func (r *userRepo) Update(u *entity.User) error {
    m := models.UserFromEntity(u)
    result := r.db.Save(m)
    if result.Error != nil {
        return fmt.Errorf("userRepo.Update: %w", result.Error)
    }
    if result.RowsAffected == 0 {
        return apperror.ErrNotFound
    }
    return nil
}