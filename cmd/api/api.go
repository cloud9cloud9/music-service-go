package api

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/zmb3/spotify"
	"music-service/internal/config"
	"music-service/internal/handler"
	"music-service/internal/repository"
	"music-service/internal/service"
	"music-service/pkg/logging"
	"net/http"
)

type Server struct {
	db     *sql.DB
	cfg    *config.Config
	client *spotify.Client
	log    *logging.LogrusLogger
}

func NewServer(
	db *sql.DB,
	cfg *config.Config,
	client *spotify.Client,
	log *logging.LogrusLogger,
) *Server {
	return &Server{
		db:     db,
		cfg:    cfg,
		client: client,
		log:    log,
	}
}

func (s *Server) Run() error {
	router := chi.NewRouter()
	repo := repository.NewRepository(s.db, s.cfg, s.log)
	services := service.NewService(repo, s.client, s.cfg.JWT.Expiration, s.cfg.JWT.Secret)
	hand := handler.NewHandler(services, s.log)
	hand.RegisterRoutes(router)
	s.log.Info("Server started on port: ", s.cfg.Server.Port)
	return http.ListenAndServe(s.cfg.Server.Port, router)
}
