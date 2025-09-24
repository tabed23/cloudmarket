package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tabed23/cloudmarket-auth/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		"localhost", "5432", "admin", "admin123", "app", "disable",
	)
	DB, err := gorm.Open(postgres.Open(dsn), initConfig())
	if err != nil {
		return nil
	}

	DB.AutoMigrate(&model.User{})
	fmt.Println("Database migrated")
	return DB
}

// InitConfig Initialize Config
func initConfig() *gorm.Config {
	return &gorm.Config{
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

// InitLog Connection Log Configuration
func initLog() logger.Interface {

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			Colorful:      true,
			LogLevel:      logger.Info,
			SlowThreshold: time.Second,
		})
	return newLogger
}
// InitNamingStrategy Init NamingStrategy
func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: false,
		TablePrefix:   "",
	}
}