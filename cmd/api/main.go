package main

import (
	"fmt"

	"ecommerce-api/config"
	delivery "ecommerce-api/internal/delivery/http"
	"ecommerce-api/internal/delivery/http/handler"
	"ecommerce-api/internal/repository/postgres"
	"ecommerce-api/internal/usecase"
)

func main() {

    cfg := config.Load()

    db := config.NewDatabase(cfg)

    productRepo := postgres.NewProductRepository(db)
    userRepo    := postgres.NewUserRepository(db)
    orderRepo   := postgres.NewOrderRepository(db)

    productUC := usecase.NewProductUsecase(productRepo)
    userUC    := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
    orderUC   := usecase.NewOrderUsecase(orderRepo, productRepo)
    _ = orderUC


    productHandler := handler.NewProductHandler(productUC)
    userHandler    := handler.NewUserHandler(userUC)

    router := delivery.NewRouter(productHandler, userHandler)
    r      := router.Setup()

    addr := ":" + cfg.ServerPort
    fmt.Printf("✓ ecommerce-api running on %s\n", addr)
    r.Run(addr)
}