package server

import "gorm.io/gorm"

func SetUpServer(db *gorm.DB) {
	go RunUserServer(db)
	go RunAuthServer(db)
}
