package handler

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"music-service/internal/models"
	"music-service/pkg/utils"
	"net/http"
	"strconv"
)

var (
	errTrackNotFound = errors.New("track not found")
	errContext       = errors.New("context error, user id not found")
)

// HandleGetTrackFromSpotify
// @Summary Get track from spotify
// @Description Get track from spotify
// @Tags tracks
// @Accept  json
// @Produce  json
// @Param trackId path string true "Track ID"
// @Success 200 {object} models.Song "Track"
// @Failure 500 {object} error "internal server error"
// @Router /tracks/{trackId} [get]
// @Security ApiKeyAuth
func (h *Handler) HandleGetTrackFromSpotify(writer http.ResponseWriter, request *http.Request) {
	trackID := chi.URLParam(request, "trackId")

	track, err := h.services.Song.GetTrackByID(trackID)
	if err != nil {
		h.log.Error("HANDLER: error getting track from spotify: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errTrackNotFound)
		return
	}

	song := utils.MapTrackToSong(track)
	h.log.Info("HANDLER: track mapped: ", song)
	utils.WriteJSON(writer, http.StatusOK, song)
}

// HandleGetTracksFromPlaylist
// @Summary Get tracks from playlist
// @Tags tracks
// @Description Get tracks from playlist
// @Accept  json
// @Produce  json
// @Param playlistId path int true "Playlist ID"
// @Success 200 {object} []models.Song "Tracks"
// @Failure 500 {object} error "internal server error"
// @Router /playlist/{playlistId}/tracks [get]
// @Security ApiKeyAuth
func (h *Handler) HandleGetTracksFromPlaylist(writer http.ResponseWriter, request *http.Request) {
	playlistId, err := strconv.Atoi(chi.URLParam(request, "playlistId"))
	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	track, err := h.services.Song.GetAllSongsFromPlaylist(userId, playlistId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(writer, "Track not found in playlist: ", http.StatusNotFound)
			return
		}
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: track founded: ", track)
	utils.WriteJSON(writer, http.StatusOK, track)
}

// HandleInsertTrackToPlaylist
// @Summary Insert track
// @Tags tracks
// @Description Insert track to playlist
// @Accept  json
// @Produce  json
// @Param playlistId path int true "Playlist ID"
// @Param trackId path string true "Track ID"
// @Success 200 {object} models.Song "Track"
// @Failure 500 {object} error "internal server error"
// @Router /tracks/{trackId}/playlist/{playlistId} [post]
// @Security ApiKeyAuth
func (h *Handler) HandleInsertTrackToPlaylist(writer http.ResponseWriter, request *http.Request) {
	playlistId, err := strconv.Atoi(chi.URLParam(request, "playlistId"))
	trackId := chi.URLParam(request, "trackId")
	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	var song models.Song

	track, err := h.services.Song.GetTrackByID(trackId)
	if err != nil {
		h.log.Error("HANDLER: error getting track from spotify: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errTrackNotFound)
		return
	}

	song = utils.MapTrackToSong(track)

	_, err = h.services.Song.CreateSong(userId, playlistId, &song)
	if err != nil {
		h.log.Error("HANDLER: error inserting track to playlist: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: track inserted to playlist: ", song)
	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"song": song,
	})
}

// HandleDeleteTrackFromPlaylist
// @Summary Delete track from playlist
// @Tags tracks
// @Description Delete track from playlist
// @Accept  json
// @Produce  json
// @Param playlistId path int true "Playlist ID"
// @Param trackId path string true "Track ID"
// @Success 200 {string} string "Track removed from playlist"
// @Failure 500 {object} error "internal server error"
// @Router /tracks/{trackId}/playlist/{playlistId} [delete]
// @Security ApiKeyAuth
func (h *Handler) HandleDeleteTrackFromPlaylist(writer http.ResponseWriter, request *http.Request) {
	playlistId, _ := strconv.Atoi(chi.URLParam(request, "playlistId"))
	trackId := chi.URLParam(request, "trackId")

	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	err = h.services.Song.DeleteSongFromPlaylist(userId, playlistId, trackId)
	if err != nil {
		h.log.Error("HANDLER: error deleting track from playlist: ", err)
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteError(writer, http.StatusInternalServerError, errTrackNotFound)
			return
		}
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: track removed from playlist")
	utils.WriteJSON(writer, http.StatusOK, "Track removed from playlist")
}
