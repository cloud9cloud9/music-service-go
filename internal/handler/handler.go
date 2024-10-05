package handler

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "music-service/docs"
	"music-service/internal/service"
	"music-service/pkg/logging"
)

const (
	apiPath              = "/api/v1"
	login                = "/login"
	register             = "/register"
	logout               = "/logout"
	ping                 = "/ping"
	playlist             = "/playlist"
	playlistById         = "/playlist/{playlistId}"
	trackFromSpotify     = "/tracks/{trackId}"
	trackFromPlayList    = "/tracks/playlist/{playlistId}"
	insertAndDeleteTrack = "/tracks/{trackId}/playlist/{playlistId}"
	swagger              = "/swagger/*"
)

type Handler struct {
	services *service.Service
	log      *logging.LogrusLogger
}

func NewHandler(services *service.Service, log *logging.LogrusLogger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route(apiPath, func(r chi.Router) {
		r.Get(swagger, httpSwagger.WrapHandler)

		r.With(h.logRequest).Post(login, h.HandleLogin)
		r.With(h.logRequest).Post(register, h.HandleRegister)
		r.With(h.userIdentity, h.logRequest).Post(logout, h.LogoutHandler)

		r.With(h.userIdentity, h.logRequest).Get(ping, h.HandlePing)

		r.With(h.userIdentity, h.logRequest).Post(playlist, h.HandleCreatePlaylist)
		r.With(h.userIdentity, h.logRequest).Get(playlist, h.HandleGetAllPlaylists)
		r.With(h.userIdentity, h.logRequest).Get(playlistById, h.HandleGetPlaylistById)
		r.With(h.userIdentity, h.logRequest).Put(playlistById, h.HandleUpdatePlaylistById)
		r.With(h.userIdentity, h.logRequest).Delete(playlistById, h.HandleDeletePlaylistById)

		r.With(h.userIdentity, h.logRequest).Get(trackFromSpotify, h.HandleGetTrackFromSpotify)
		r.With(h.userIdentity, h.logRequest).Get(trackFromPlayList, h.HandleGetTracksFromPlaylist)
		r.With(h.userIdentity, h.logRequest).Post(insertAndDeleteTrack, h.HandleInsertTrackToPlaylist)
		r.With(h.userIdentity, h.logRequest).Delete(insertAndDeleteTrack, h.HandleDeleteTrackFromPlaylist)

	})
}
