package handler

import (
	"errors"
	"net/http"
	"strconv"

	"ecommerce-api/internal/domain/entity"
	delivery "ecommerce-api/internal/delivery/http/responses"
	"ecommerce-api/pkg/apperror"

	"github.com/gin-gonic/gin"
)

func mapErr(c *gin.Context, err error) {
    switch {
    case errors.Is(err, apperror.ErrNotFound):
        delivery.Fail(c, http.StatusNotFound, err.Error())
    case errors.Is(err, apperror.ErrUnauthorized):
        delivery.Fail(c, http.StatusUnauthorized, err.Error())
    case errors.Is(err, apperror.ErrConflict):
        delivery.Fail(c, http.StatusConflict, err.Error())
    case errors.Is(err, apperror.ErrOutOfStock):
        delivery.Fail(c, http.StatusConflict, err.Error())
    case errors.Is(err, apperror.ErrInvalidInput),
         errors.Is(err, apperror.ErrBadRequest):
        delivery.Fail(c, http.StatusBadRequest, err.Error())
    default:
        delivery.Fail(c, http.StatusInternalServerError, "internal server error")
    }
}

func parseID(c *gin.Context) (uint, bool) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        delivery.Fail(c, http.StatusBadRequest, "invalid id — must be a positive integer")
        return 0, false
    }
    return uint(id), true
}

type createProductReq struct {
    Name        string  `json:"name"        binding:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price"       binding:"required,gt=0"`
    Stock       int     `json:"stock"       binding:"min=0"`
    CategoryID  uint    `json:"category_id"`
}

type updateProductReq struct {
    Name        string  `json:"name"        binding:"required"`
    Description string  `json:"description"`
    Price       float64 `json:"price"       binding:"required,gt=0"`
    Stock       int     `json:"stock"       binding:"min=0"`
    CategoryID  uint    `json:"category_id"`
}

type ProductHandler struct {
    uc entity.ProductUsecase
}

func NewProductHandler(uc entity.ProductUsecase) *ProductHandler {
    return &ProductHandler{uc: uc}
}

// GET /api/v1/products
func (h *ProductHandler) List(c *gin.Context) {
    page,  _ := strconv.Atoi(c.DefaultQuery("page",  "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    products, total, err := h.uc.ListProducts(page, limit)
    if err != nil { mapErr(c, err); return }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    products,
        "total":   total,
        "page":    page,
        "limit":   limit,
    })
}

// GET /api/v1/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    product, err := h.uc.GetProduct(id)
    if err != nil { mapErr(c, err); return }

    delivery.OK(c, product)
}

// POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
    var req createProductReq
    if err := c.ShouldBindJSON(&req); err != nil {
        delivery.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    product, err := h.uc.CreateProduct(entity.CreateProductRequest{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        CategoryID:  req.CategoryID,
    })
    if err != nil { mapErr(c, err); return }

    delivery.Created(c, product)
}

// PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    var req updateProductReq
    if err := c.ShouldBindJSON(&req); err != nil {
        delivery.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    product, err := h.uc.UpdateProduct(id, entity.UpdateProductRequest{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        CategoryID:  req.CategoryID,
    })
    if err != nil { mapErr(c, err); return }

    delivery.OK(c, product)
}

// DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    if err := h.uc.DeleteProduct(id); err != nil {
        mapErr(c, err)
        return
    }
    delivery.OK(c, gin.H{"message": "product deleted"})
}