package responses

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "strings"
    "github.com/go-playground/validator/v10"
    "fmt"
    "errors"

)

type Response struct {
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
    Error   string `json:"error,omitempty"`
    Message string `json:"message,omitempty"`
}

type PaginatedResponse struct {
    Success bool  `json:"success"`
    Data    any   `json:"data"`
    Meta    Meta  `json:"meta"`
}

type Meta struct {
    Total      int64 `json:"total"`
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    TotalPages int   `json:"total_pages"`
}

func OK(c *gin.Context, data any) {
    c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func Created(c *gin.Context, data any) {
    c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

func Message(c *gin.Context, msg string) {
    c.JSON(http.StatusOK, Response{Success: true, Message: msg})
}

func Paginated(c *gin.Context, data any, total int64, page, limit int) {
    totalPages := int(total) / limit
    if int(total)%limit != 0 { totalPages++ }
    c.JSON(http.StatusOK, PaginatedResponse{
        Success: true,
        Data:    data,
        Meta:    Meta{Total: total, Page: page, Limit: limit, TotalPages: totalPages},
    })
}

func Fail(c *gin.Context, status int, msg string) {
    c.JSON(status, Response{Success: false, Error: msg})
}


type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func FailValidation(c *gin.Context, err error) {
    var ve validator.ValidationErrors
    if !errors.As(err, &ve) {
        Fail(c, http.StatusBadRequest, err.Error())
        return
    }

    errs := make([]ValidationError, len(ve))
    for i, fe := range ve {
        errs[i] = ValidationError{
            Field:   strings.ToLower(fe.Field()),
            Message: validationMsg(fe),
        }
    }
    c.JSON(http.StatusBadRequest, gin.H{
        "success": false,
        "error":   "validation failed",
        "details": errs,
    })
}

func validationMsg(fe validator.FieldError) string {
    switch fe.Tag() {
    case "required": return "this field is required"
    case "email":    return "must be a valid email address"
    case "min":      return fmt.Sprintf("minimum value is %s", fe.Param())
    case "max":      return fmt.Sprintf("maximum value is %s", fe.Param())
    case "gt":       return fmt.Sprintf("must be greater than %s", fe.Param())
    case "gte":      return fmt.Sprintf("must be at least %s", fe.Param())
    default:          return fmt.Sprintf("failed validation: %s", fe.Tag())
    }
}