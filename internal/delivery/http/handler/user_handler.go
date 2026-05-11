package handler

import (
    "net/http"

    "ecommerce-api/internal/domain/contract"
    delivery "ecommerce-api/internal/delivery/http/responses"

    "github.com/gin-gonic/gin"
)

type registerReq struct {
    Name     string `json:"name"     binding:"required"`
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type loginReq struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type UserHandler struct {
    uc contract.UserUsecase
}

func NewUserHandler(uc contract.UserUsecase) *UserHandler {
    return &UserHandler{uc: uc}
}

// POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
    var req registerReq
    if err := c.ShouldBindJSON(&req); err != nil {
        delivery.Fail(c, http.StatusBadRequest, err.Error())
        return
    }
    user, err := h.uc.Register(contract.RegisterRequest{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil { mapErr(c, err); return }
    delivery.Created(c, gin.H{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "role":  user.Role,
    })
}

// POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
    var req loginReq
    if err := c.ShouldBindJSON(&req); err != nil {
        delivery.Fail(c, http.StatusBadRequest, err.Error())
        return
    }
    token, err := h.uc.Login(contract.LoginRequest{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil { mapErr(c, err); return }
    delivery.OK(c, gin.H{"token": token})
}