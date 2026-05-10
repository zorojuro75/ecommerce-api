package main

import (
    "fmt"
    "log"

    "ecommerce-api/config"
    delivery "ecommerce-api/internal/delivery/http"
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

    _ = productUC
    _ = userUC
    _ = orderUC

    r := delivery.SetupRouter()

    fmt.Printf("✓ ecommerce-api running on :%s\n", cfg.ServerPort)
    log.Fatal(r.Run(":" + cfg.ServerPort))
}