package handler

import (
	"github.com/go-chi/chi/v5"
	"music-service/internal/models"
	"music-service/pkg/utils"
	"net/http"
	"strconv"
)

// HandleCreatePlaylist
//
// @Summary Create playlist
// @Description Creates a new playlist for the authenticated user
// @Tags playlist
// @Accept  json
// @Produce  json
// @Param input body models.CreatePlaylistDto true "Playlist creation dto"
// @Success 200 {object} map[string]interface{} "Playlist created"
// @Failure 400 {object} error "invalid parsing JSON"
// @Failure 500 {object} error "internal server error"
// @Router /playlist [post]
// @Security ApiKeyAuth
func (h *Handler) HandleCreatePlaylist(writer http.ResponseWriter, request *http.Request) {
	var input models.CreatePlaylistDto
	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	if err := utils.ParseJSON(request, &input); err != nil {
		h.log.Error("HANDLER: error parsing JSON: ", err)
		utils.InvalidParsingJSON(writer)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		h.log.Error("HANDLER: error validating input: ", err)
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	playlist := &models.Playlist{
		Name:   input.Name,
		UserId: userId,
	}

	id, err := h.services.PlayList.CreatePlaylist(playlist)
	if err != nil {
		h.log.Error("HANDLER: error creating playlist: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: playlist created: ", id)
	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"id":     id,
	})
}

// HandleGetAllPlaylists
// @Summary Get all playlists
// @Tags playlist
// @Description Get all playlists for the authenticated user
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Playlist "Playlists"
// @Failure 500 {object} error "internal server error"
// @Router /playlist [get]
// @Security ApiKeyAuth
func (h *Handler) HandleGetAllPlaylists(writer http.ResponseWriter, request *http.Request) {
	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	playlists, err := h.services.PlayList.GetAllPlaylists(userId)
	if err != nil {
		h.log.Error("HANDLER: error getting playlists: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: playlists found: ", playlists)
	utils.WriteJSON(writer, http.StatusOK, playlists)
}

// HandleGetPlaylistById
// @Summary Get playlist by id
// @Tags playlist
// @Description Get playlist by id
// @Accept  json
// @Produce  json
// @Param id path int true "Playlist id"
// @Success 200 {object} models.Playlist "Playlist"
// @Failure 500 {object} error "internal server error"
// @Router /playlist/{id} [get]
// @Security ApiKeyAuth
func (h *Handler) HandleGetPlaylistById(writer http.ResponseWriter, request *http.Request) {
	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	playlistId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		h.log.Error("HANDLER: error getting playlist id: ", err)
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	playlist, err := h.services.PlayList.GetPlaylistById(userId, playlistId)
	if err != nil {
		h.log.Error("HANDLER: error getting playlist from db: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: playlist found: ", playlist)
	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"playlist": playlist,
	})
}

// HandleUpdatePlaylistById
// @Summary Update playlist by id
// @Tags playlist
// @Description Update playlist by id
// @Accept  json
// @Produce  json
// @Param id path int true "Playlist id"
// @Param input body models.UpdatePlaylistDto true "Playlist update dto"
// @Success 200 {object} map[string]interface{} "Playlist updated"
// @Failure 400 {object} error "invalid parsing JSON"
// @Failure 500 {object} error "internal server error"
// @Router /playlist/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) HandleUpdatePlaylistById(writer http.ResponseWriter, request *http.Request) {
	var input models.UpdatePlaylistDto

	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	if err := utils.ParseJSON(request, &input); err != nil {
		h.log.Error("HANDLER: error parsing JSON: ", err)
		utils.InvalidParsingJSON(writer)
		return
	}

	playlistId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		h.log.Error("HANDLER: error getting playlist id: ", err)
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	updated := &models.Playlist{
		Name: input.Name,
		ID:   playlistId,
	}

	err = h.services.PlayList.UpdatePlaylistById(userId, updated)
	if err != nil {
		h.log.Error("HANDLER: error updating playlist: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: playlist updated: ", updated)
	utils.WriteJSON(writer, http.StatusOK, updated)
}

// HandleDeletePlaylistById
// @Summary Delete playlist by id
// @Tags playlist
// @Description Delete playlist by id
// @Accept  json
// @Produce  json
// @Param id path int true "Playlist id"
// @Success 200 {object} map[string]interface{} "Playlist deleted"
// @Failure 400 {object} error "invalid playlist id"
// @Failure 500 {object} error "internal server error"
// @Router /playlist/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) HandleDeletePlaylistById(writer http.ResponseWriter, request *http.Request) {
	userId, err := getUserId(request.Context())
	if err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errContext)
		return
	}

	playlistId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		h.log.Error("HANDLER: error getting playlist id: ", err)
		utils.WriteError(writer, http.StatusBadRequest, err)
		return
	}

	err = h.services.PlayList.DeletePlaylistById(userId, playlistId)
	if err != nil {
		h.log.Error("HANDLER: error deleting playlist: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, err)
		return
	}

	h.log.Info("HANDLER: playlist deleted: ", playlistId)
	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"message": "playlist deleted",
	})
}
