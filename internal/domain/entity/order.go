package entity

import (
    "time"
    "ecommerce-api/pkg/apperror"
)

type OrderStatus string

const (
    StatusPending   OrderStatus = "pending"
    StatusConfirmed OrderStatus = "confirmed"
    StatusShipped   OrderStatus = "shipped"
    StatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
    ProductID   uint
    ProductName string
    Price       float64
    Quantity    int
    Subtotal    float64
}

type Order struct {
    ID        uint
    UserID    uint
    Items     []OrderItem
    Total     float64
    Status    OrderStatus
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (o *Order) CalculateTotal() {
    total := 0.0
    for _, item := range o.Items {
        total += item.Price * float64(item.Quantity)
    }
    o.Total = total
}

func (o Order) CanCancel() bool {
    return o.Status == StatusPending
}

func (o Order) Validate() error {
    if o.UserID == 0       { return apperror.ErrInvalidInput }
    if len(o.Items) == 0  { return apperror.ErrInvalidInput }
    return nil
}

