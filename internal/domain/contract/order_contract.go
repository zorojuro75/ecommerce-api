package contract 

import(
	"ecommerce-api/internal/domain/entity"
)

type OrderUsecase interface {
    PlaceOrder(req PlaceOrderRequest) (*entity.Order, error)
    GetOrder(id uint) (*entity.Order, error)
    ListUserOrders(userID uint, page, limit int) ([]entity.Order, int64, error)
    CancelOrder(id, userID uint) error
}

type PlaceOrderRequest struct {
    UserID uint
    Items  []OrderItemRequest
}

type OrderItemRequest struct {
    ProductID uint
    Quantity  int
}