package database

import (
	"fmt"
	"os"

	"github.com/BraunKc/todo/db-service/config"
	"github.com/BraunKc/todo/db-service/internal/repository/models"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		config.Logger.Fatal("no .env file found", zap.Error(err))
	}

	host, user, password, dbname, port := os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		config.Logger.Fatal("DB open error", zap.Error(err))
	}

	db.AutoMigrate(&models.Task{})

	config.Logger.Debug("database inited", zap.Any("DB", db))
	return db
}
