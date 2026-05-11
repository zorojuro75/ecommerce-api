package usecase

import (
    "fmt"

    "ecommerce-api/internal/domain/entity"
    "ecommerce-api/internal/domain/contract"
    domainrepo "ecommerce-api/internal/domain/repository"
    "ecommerce-api/pkg/apperror"
)

type orderUsecase struct {
    orderRepo   domainrepo.OrderRepository
    productRepo domainrepo.ProductRepository
}

func NewOrderUsecase(
    orderRepo   domainrepo.OrderRepository,
    productRepo domainrepo.ProductRepository,
) contract.OrderUsecase {
    return &orderUsecase{
        orderRepo:   orderRepo,
        productRepo: productRepo,
    }
}

func (uc *orderUsecase) PlaceOrder(req contract.PlaceOrderRequest) (*entity.Order, error) {
    if len(req.Items) == 0 {
        return nil, fmt.Errorf("PlaceOrder: %w", apperror.ErrInvalidInput)
    }

    orderItems := make([]entity.OrderItem, 0, len(req.Items))

    for _, itemReq := range req.Items {
        if itemReq.Quantity <= 0 {
            return nil, fmt.Errorf("PlaceOrder: quantity must be > 0: %w", apperror.ErrInvalidInput)
        }

        product, err := uc.productRepo.FindByID(itemReq.ProductID)
        if err != nil {
            return nil, fmt.Errorf("PlaceOrder: product %d: %w", itemReq.ProductID, err)
        }

        if product.Stock < itemReq.Quantity {
            return nil, fmt.Errorf("PlaceOrder: %s: %w", product.Name, apperror.ErrOutOfStock)
        }

        orderItems = append(orderItems, entity.OrderItem{
            ProductID:   product.ID,
            ProductName: product.Name,
            Price:       product.Price,
            Quantity:    itemReq.Quantity,
            Subtotal:    product.Price * float64(itemReq.Quantity),
        })

        product.Stock -= itemReq.Quantity
        if err := uc.productRepo.Update(product); err != nil {
            return nil, fmt.Errorf("PlaceOrder: update stock: %w", err)
        }
    }

    order := &entity.Order{
        UserID: req.UserID,
        Items:  orderItems,
        Status: entity.StatusPending,
    }
    order.CalculateTotal()

    if err := order.Validate(); err != nil {
        return nil, fmt.Errorf("PlaceOrder: %w", err)
    }
    if err := uc.orderRepo.Create(order); err != nil {
        return nil, fmt.Errorf("PlaceOrder: %w", err)
    }
    return order, nil
}

func (uc *orderUsecase) GetOrder(id uint) (*entity.Order, error) {
    o, err := uc.orderRepo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetOrder id=%d: %w", id, err)
    }
    return o, nil
}

func (uc *orderUsecase) ListUserOrders(userID uint, page, limit int) ([]entity.Order, int64, error) {
    if page < 1  { page  = 1 }
    if limit < 1 { limit = 10 }
    orders, total, err := uc.orderRepo.FindByUserID(userID, page, limit)
    if err != nil {
        return nil, 0, fmt.Errorf("ListUserOrders: %w", err)
    }
    return orders, total, nil
}

func (uc *orderUsecase) CancelOrder(id, userID uint) error {
    o, err := uc.orderRepo.FindByID(id)
    if err != nil {
        return fmt.Errorf("CancelOrder: %w", err)
    }

    if o.UserID != userID {
        return fmt.Errorf("CancelOrder: %w", apperror.ErrUnauthorized)
    }

    if !o.CanCancel() {
        return fmt.Errorf("CancelOrder: order is %s, cannot cancel: %w",
            o.Status, apperror.ErrBadRequest)
    }
    return uc.orderRepo.UpdateStatus(id, entity.StatusCancelled)
}