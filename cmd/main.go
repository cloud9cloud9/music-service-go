package main

import (
	"github.com/go-sql-driver/mysql"
	"music-service/cmd/api"
	"music-service/internal/config"
	"music-service/pkg/client"
	"music-service/pkg/db"
	"music-service/pkg/logging"

	_ "github.com/go-sql-driver/mysql"
)

// @title Music API
// @version 1.0
// @description This is a sample server for managing users/playlists and songs.
// @host localhost:8082
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()
	log := logging.NewLogger()
	svc := client.NewSpotifyClient(cfg.Spotify.ClientID, cfg.Spotify.ClientSecret)
	storage, err := db.NewStorage(mysql.Config{
		User:                 cfg.MySql.DBUser,
		Passwd:               cfg.MySql.DBPassword,
		Net:                  cfg.MySql.DBNet,
		Addr:                 cfg.MySql.DBAddress,
		DBName:               cfg.MySql.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	server := api.NewServer(storage, cfg, svc, log)
	if err != nil {
		log.Error("Database error: ", err)
	}
	log.Info("Database connected")
	if err := server.Run(); err != nil {
		log.Error("Server error: %s", err)
	}
}
