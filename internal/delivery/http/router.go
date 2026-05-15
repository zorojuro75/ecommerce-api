package http

import (
	"fmt"
	"net/http"
	"time"

	"ecommerce-api/internal/delivery/http/handler"
	"ecommerce-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
    productHandler *handler.ProductHandler
    userHandler    *handler.UserHandler
    orderHandler   *handler.OrderHandler
    jwtSecret      string
}

func NewRouter(
    productHandler *handler.ProductHandler,
    userHandler    *handler.UserHandler,
    orderHandler   *handler.OrderHandler,
    jwtSecret      string,
) *Router {
    return &Router{
        productHandler: productHandler,
        userHandler:    userHandler,
        orderHandler:   orderHandler,
        jwtSecret:      jwtSecret,
    }
}

func (ro *Router) Setup() *gin.Engine {
    r := gin.New()
    r.Use(middleware.Recovery())
    r.Use(middleware.Logger())
    r.NoRoute(func(c *gin.Context){
        c.JSON(404, gin.H{
            "success": false,
            "error": fmt.Sprintf("route %s %s not found", c.Request.Method, c.Request.URL.Path),
        })
    })
    r.NoMethod(func(c *gin.Context) {
        c.JSON(405, gin.H{
            "success": false,
            "error":   "method not allowed",
        })
    })

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

    auth := middleware.Auth(ro.jwtSecret)
    protected := v1.Group("", auth)
    {
        p := protected.Group("/products")
        p.POST("",       ro.productHandler.Create)
        p.PUT(("/:id"),  ro.productHandler.Update)
        p.DELETE("/:id", ro.productHandler.Delete)

        v1.GET("/me", auth, ro.userHandler.Me)
    }
    orders := protected.Group("/orders")
    {
        orders.POST("",            ro.orderHandler.PlaceOrder)
        orders.GET("/my",         ro.orderHandler.MyOrders)
        orders.GET("/:id",        ro.orderHandler.GetOrder)
        orders.PATCH("/:id/cancel", ro.orderHandler.CancelOrder)
    }
    return r
}