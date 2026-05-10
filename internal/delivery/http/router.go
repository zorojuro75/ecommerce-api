package http

import (
	"net/http"
	"time"

	"ecommerce-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
    productHandler *handler.ProductHandler
    userHandler    *handler.UserHandler
}

func NewRouter(
    productHandler *handler.ProductHandler,
    userHandler    *handler.UserHandler,
) *Router {
    return &Router{
        productHandler: productHandler,
        userHandler:    userHandler,
    }
}

func (ro *Router) Setup() *gin.Engine {
    r := gin.Default()

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
        {
            auth.POST("/register", ro.userHandler.Register)
            auth.POST("/login",    ro.userHandler.Login)
        }

        products := v1.Group("/products")
        {
            products.GET("",     ro.productHandler.List)
            products.GET("/:id", ro.productHandler.Get)
            products.POST("",    ro.productHandler.Create)
            products.PUT("/:id", ro.productHandler.Update)
            products.DELETE("/:id", ro.productHandler.Delete)
        }
    }

    return r
}

func SetupRouter() *gin.Engine {
    return gin.Default()
}