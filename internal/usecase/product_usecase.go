package usecase

import (
    "fmt"

    "ecommerce-api/internal/domain/entity"
    "ecommerce-api/internal/domain/contract"
    domainrepo "ecommerce-api/internal/domain/repository"
)

type productUsecase struct {
    repo domainrepo.ProductRepository
}

func NewProductUsecase(repo domainrepo.ProductRepository) contract.ProductUsecase {
    return &productUsecase{repo: repo}
}

func (uc *productUsecase) CreateProduct(req contract.CreateProductRequest) (*entity.Product, error) {
    p := &entity.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        CategoryID:  req.CategoryID,
    }
    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("CreateProduct: %w", err)
    }
    if err := uc.repo.Create(p); err != nil {
        return nil, fmt.Errorf("CreateProduct: %w", err)
    }
    return p, nil
}

func (uc *productUsecase) GetProduct(id uint) (*entity.Product, error) {
    p, err := uc.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetProduct id=%d: %w", id, err)
    }
    return p, nil
}

func (uc *productUsecase) ListProducts(filter entity.ProductFilter) ([]entity.Product, int64, error) {
    if filter.Page  < 1   { filter.Page  = 1 }
    if filter.Limit < 1   { filter.Limit = 10 }
    if filter.Limit > 100 { filter.Limit = 100 }

    products, total, err := uc.repo.FindAll(filter)
    if err != nil {
        return nil, 0, fmt.Errorf("ListProducts: %w", err)
    }
    return products, total, nil
}

func (uc *productUsecase) UpdateProduct(id uint, req contract.UpdateProductRequest) (*entity.Product, error) {
    p, err := uc.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("UpdateProduct: %w", err)
    }
    p.Name        = req.Name
    p.Description = req.Description
    p.Price       = req.Price
    p.Stock       = req.Stock
    p.CategoryID  = req.CategoryID

    if err := p.Validate(); err != nil {
        return nil, fmt.Errorf("UpdateProduct: %w", err)
    }
    if err := uc.repo.Update(p); err != nil {
        return nil, fmt.Errorf("UpdateProduct: %w", err)
    }
    return p, nil
}

func (uc *productUsecase) DeleteProduct(id uint) error {
    if _, err := uc.repo.FindByID(id); err != nil {
        return fmt.Errorf("DeleteProduct: %w", err)
    }
    return uc.repo.Delete(id)
}