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

func Auth(jwtSecret string) gin.HandlerFunc {
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
                "error":   "authorization format must be: Bearer <token>",
            })
            return
        }

        tokenStr := parts[1]

        claims, err := jwt.Validate(tokenStr, jwtSecret)
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

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, _ := c.Get(ContextRole)
        if role != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "success": false,
                "error":   "admin access required",
            })
            return
        }
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

func GetRole(c *gin.Context) string {
    val, _ := c.Get(ContextRole)
    role, _ := val.(string)
    return role
}