package db

import (
	"adiubaidah/adi-bot/helper"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() *gorm.DB {
	dbHost := helper.GetEnv("DB_HOST")
	dbUsername := helper.GetEnv("DB_USER")
	dbPassword := helper.GetEnv("DB_PASS")
	dbName := helper.GetEnv("DB_NAME")
	dbPort := helper.GetEnv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Jakarta", dbHost, dbUsername, dbPassword, dbName, dbPort)
	dialect := postgres.New(postgres.Config{
		DSN: dsn,
		// PreferSimpleProtocol: ,
	})
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal(err)
	}

	// Get the underlying sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db
}
