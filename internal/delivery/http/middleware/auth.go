package middleware

import (
    "net/http"
    "strings"

    "ecommerce-api/pkg/jwt"

    "github.com/gin-gonic/gin"
)

const (
    ContextUserID = "userID"
    ContextRole   = "role"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "authorization header required",
            })
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "authorization format: Bearer <token>",
            })
            return
        }
        tokenStr := parts[1]

        claims, err := jwt.Validate(tokenStr, secret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "invalid or expired token",
            })
            return
        }

        c.Set(ContextUserID, claims.UserID)
        c.Set(ContextRole,   claims.Role)

        c.Next()
    }
}

func GetUserID(c *gin.Context) (uint, bool) {
    val, exists := c.Get(ContextUserID)
    if !exists {
        return 0, false
    }
    id, ok := val.(uint)
    return id, ok
}

func GetRole(c *gin.Context) (string, bool) {
    val, exists := c.Get(ContextRole)
    if !exists {
        return "", false
    }
    role, ok := val.(string)
    return role, ok
}

func MustGetUserID(c *gin.Context) uint {
    id, ok := GetUserID(c)
    if !ok {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "success": false,
            "error":   "unauthorized",
        })
        return 0
    }
    return id
}