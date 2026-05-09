package models

import (
    "ecommerce-api/internal/domain/entity"
    "gorm.io/gorm"
)


type OrderItemModel struct {
	gorm.Model
	OrderID     		uint    	`gorm:"not null;index"`
	ProductID			uint		`gorm:"notnull"`
	ProductName			string		`gorm:"not null"`
	Price				float64		`gorm:"not null"`
	Quantity			int			`gorm:"not null"`
	Subtotal			float64		`gorm:"not null"`
}


func (OrderItemModel) TableName() string { return "order_items"}

func (m *OrderItemModel) OderItemModelToEntity() *entity.OrderItem{
	return &entity.OrderItem{
		ProductID: m.ProductID,
		ProductName: m.ProductName,
		Price: m.Price,
		Quantity: m.Quantity,
		Subtotal: m.Subtotal,
	}
}

func EntityToOrderItemModel(o *entity.OrderItem) *OrderItemModel{
	return &OrderItemModel{
		Model: gorm.Model{ID: o.ProductID},
		ProductName: o.ProductName,
		Price: o.Price,
		Quantity: o.Quantity,
		Subtotal: o.Subtotal,
	}
}

type OrderModel struct {
	gorm.Model
	UserID 		uint             `gorm:"not null;index"`
	Total  		float64          `gorm:"not null"`
	Status 		string           `gorm:"type:varchar(20);not null;default:'pending'"`
	Items  		[]OrderItemModel `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (OrderModel) TableName() string { return "orders" }

func (m *OrderModel) OrderModelToEntity() *entity.Order {
	items := make([]entity.OrderItem, len(m.Items))
	for i, item := range m.Items {
		items[i] = *item.OderItemModelToEntity()
	}

	return &entity.Order{
		ID:        m.ID,
		UserID:    m.UserID,
		Total:     m.Total,
		Status:    entity.OrderStatus(m.Status),
		Items:     items,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func EntityToOrderModel(o *entity.Order) *OrderModel {
	items := make([]OrderItemModel, len(o.Items))
	for i, item := range o.Items {
		modelItem := EntityToOrderItemModel(&item)
		modelItem.OrderID = o.ID
		items[i] = *modelItem
	}

	return &OrderModel{
		Model:  gorm.Model{ID: o.ID},
		UserID: o.UserID,
		Total:  o.Total,
		Status: string(o.Status),
		Items:  items,
	}
}