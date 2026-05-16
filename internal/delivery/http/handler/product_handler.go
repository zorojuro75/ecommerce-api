package handler

import (
	"errors"
	"net/http"
	"strconv"

	responses "ecommerce-api/internal/delivery/http/responses"
	"ecommerce-api/internal/domain/contract"
	"ecommerce-api/internal/domain/entity"
	"ecommerce-api/pkg/apperror"
	"ecommerce-api/pkg/pagination"

	"github.com/gin-gonic/gin"
)

func mapErr(c *gin.Context, err error) {
    switch {
    case errors.Is(err, apperror.ErrNotFound):
        responses.Fail(c, http.StatusNotFound, err.Error())
    case errors.Is(err, apperror.ErrUnauthorized):
        responses.Fail(c, http.StatusUnauthorized, err.Error())
    case errors.Is(err, apperror.ErrConflict):
        responses.Fail(c, http.StatusConflict, err.Error())
    case errors.Is(err, apperror.ErrOutOfStock):
        responses.Fail(c, http.StatusConflict, err.Error())
    case errors.Is(err, apperror.ErrInvalidInput),
         errors.Is(err, apperror.ErrBadRequest):
         responses.Fail(c, http.StatusBadRequest, err.Error())
    default:
        responses.Fail(c, http.StatusInternalServerError, "internal server error")
    }
}

func parseID(c *gin.Context) (uint, bool) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        responses.Fail(c, http.StatusBadRequest, "invalid id — must be a positive integer")
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
    uc contract.ProductUsecase
}

func NewProductHandler(uc contract.ProductUsecase) *ProductHandler {
    return &ProductHandler{uc: uc}
}

// GET /api/v1/products
func (h *ProductHandler) List(c *gin.Context) {
    p := pagination.FromContext(c)

    minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
    maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

    filter := entity.ProductFilter{
        Page:     p.Page,
        Limit:    p.Limit,
        Search:   c.Query("search"),
        MinPrice: minPrice,
        MaxPrice: maxPrice,
        Sort:     c.Query("sort"),
    }

    products, total, err := h.uc.ListProducts(filter)
    if err != nil { mapErr(c, err); return }

    responses.Paginated(c, products, total, p.Page, p.Limit)
}

// GET /api/v1/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    product, err := h.uc.GetProduct(id)
    if err != nil { mapErr(c, err); return }

    responses.OK(c, product)
}

// POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
    var req createProductReq
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    product, err := h.uc.CreateProduct(contract.CreateProductRequest{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        CategoryID:  req.CategoryID,
    })
    if err != nil { mapErr(c, err); return }

    responses.Created(c, product)
}

// PUT /api/v1/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    var req updateProductReq
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    product, err := h.uc.UpdateProduct(id, contract.UpdateProductRequest{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        CategoryID:  req.CategoryID,
    })
    if err != nil { mapErr(c, err); return }

    responses.OK(c, product)
}

// DELETE /api/v1/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
    id, ok := parseID(c)
    if !ok { return }

    if err := h.uc.DeleteProduct(id); err != nil {
        mapErr(c, err)
        return
    }
    responses.OK(c, gin.H{"message": "product deleted"})
}