package config

import (
    "fmt"
    "log"

    "ecommerce-api/internal/repository/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func NewDatabase(cfg *Config) *gorm.DB {
    if cfg.DatabaseURL == "" {
        log.Fatal("DATABASE_URL is required")
    }

    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }

    // AutoMigrate — creates/updates tables from models
    err = db.AutoMigrate(
        &models.ProductModel{},
        &models.UserModel{},
        &models.OrderModel{},
        &models.OrderItemModel{},
    )
    if err != nil {
        log.Fatalf("migration failed: %v", err)
    }

    fmt.Println("✓ Database connected and migrated")
    return db
}