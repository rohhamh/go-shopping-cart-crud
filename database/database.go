package database

import (
	"fmt"

	"github.com/rohhamh/go-shopping-cart-crud/config"
	"github.com/rohhamh/go-shopping-cart-crud/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)
var DB *gorm.DB

func Connect()  *gorm.DB {
	sqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password,
	)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Cart{})
	if err != nil {
		panic(err)
	}

	DB = db
	return DB
}

func Connection() *gorm.DB {
	return DB
}
