package main

import (
	"grpcapp/db"
	models "grpcapp/model"
	"grpcapp/server"
	"grpcapp/utils"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load .env")
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	db, err := db.InitDb(config)

	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}
	// Migrate the database schema
	models.Migrate(db)
	server.SetUpServer(db)
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
	<-gracefulShutdown
}
