package http

import (
    "net/http"
    "time"

    "ecommerce-api/internal/delivery/http/handler"
    "ecommerce-api/internal/delivery/http/middleware"

    "github.com/gin-gonic/gin"
)

type Router struct {
    productHandler *handler.ProductHandler
    userHandler    *handler.UserHandler
    jwtSecret      string
}

func NewRouter(
    productHandler *handler.ProductHandler,
    userHandler    *handler.UserHandler,
    jwtSecret      string,
) *Router {
    return &Router{
        productHandler: productHandler,
        userHandler:    userHandler,
        jwtSecret:      jwtSecret,
    }
}

func (ro *Router) Setup() *gin.Engine {
    r := gin.Default()

    // Health — no auth
    startTime := time.Now()
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "status":  "ok",
            "uptime":  time.Since(startTime).String(),
        })
    })

    v1 := r.Group("/api/v1")

    {
        auth := v1.Group("/auth")
        auth.POST("/register", ro.userHandler.Register)
        auth.POST("/login",    ro.userHandler.Login)

        products := v1.Group("/products")
        products.GET("",     ro.productHandler.List)
        products.GET("/:id", ro.productHandler.Get)
    }

    auth := middleware.AuthMiddleware(ro.jwtSecret)
    protected := v1.Group("", auth)
    {
        p := protected.Group("/products")
        p.POST("",       ro.productHandler.Create)
        p.PUT(("/:id"),  ro.productHandler.Update)
        p.DELETE("/:id", ro.productHandler.Delete)

        v1.GET("/me", auth, ro.userHandler.Me)
    }

    return r
}