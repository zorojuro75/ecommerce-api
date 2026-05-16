package handler

import (
	"net/http"

	"ecommerce-api/internal/delivery/http/middleware"
	responses "ecommerce-api/internal/delivery/http/responses"
	"ecommerce-api/internal/domain/contract"
	"ecommerce-api/pkg/pagination"

	"github.com/gin-gonic/gin"
)

type orderItemReq struct {
    ProductID uint `json:"product_id" binding:"required"`
    Quantity  int  `json:"quantity"   binding:"required,min=1"`
}

type placeOrderReq struct {
    Items []orderItemReq `json:"items" binding:"required,min=1,dive"`
}

// Handler struct
type OrderHandler struct {
    uc contract.OrderUsecase
}

func NewOrderHandler(uc contract.OrderUsecase) *OrderHandler {
    return &OrderHandler{uc: uc}
}

// POST /api/v1/orders
func (h *OrderHandler) PlaceOrder(c *gin.Context) {
    userID, ok := middleware.GetUserID(c)
    if !ok {
        responses.Fail(c, http.StatusUnauthorized, "unauthorized")
        return
    }

    var req placeOrderReq
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    items := make([]contract.OrderItemRequest, len(req.Items))
    for i, item := range req.Items {
        items[i] = contract.OrderItemRequest{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
        }
    }

    order, err := h.uc.PlaceOrder(contract.PlaceOrderRequest{
        UserID: userID,
        Items:  items,
    })

    if err != nil { mapErr(c, err); return }

    responses.Created(c, order)
}

// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
    userID, ok := middleware.GetUserID(c)
    if !ok {
        responses.Fail(c, http.StatusUnauthorized, "unauthorized")
        return
    }

    id, ok := parseID(c)
    if !ok { return }

    order, err := h.uc.GetOrder(id)
    if err != nil { mapErr(c, err); return }

    if order.UserID != userID {
        responses.Fail(c, http.StatusForbidden, "access denied")
        return
    }

    responses.OK(c, order)
}

// GET /api/v1/orders/my?page=1&limit=10
func (h *OrderHandler) MyOrders(c *gin.Context) {
    userID, ok := middleware.GetUserID(c)
    if !ok {
        responses.Fail(c, http.StatusUnauthorized, "unauthorized")
        return
    }

    p := pagination.FromContext(c)

    orders, total, err := h.uc.ListUserOrders(userID, p.Page, p.Limit)
    if err != nil { mapErr(c, err); return }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    orders,
        "total":   total,
        "page":    p.Page,
        "limit":   p.Limit,
    })
}

// PATCH /api/v1/orders/:id/cancel
func (h *OrderHandler) CancelOrder(c *gin.Context) {
    userID, ok := middleware.GetUserID(c)
    if !ok {
        responses.Fail(c, http.StatusUnauthorized, "unauthorized")
        return
    }

    id, ok := parseID(c)
    if !ok { return }

    if err := h.uc.CancelOrder(id, userID); err != nil {
        mapErr(c, err)
        return
    }

    responses.OK(c, gin.H{"message": "order cancelled"})
}