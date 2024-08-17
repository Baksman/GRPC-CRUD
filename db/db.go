package db

import (
	"grpcapp/utils"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(config utils.Config) (db *gorm.DB, err error) {
	var dsn string = config.DBSource
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}
	return db, err
}
