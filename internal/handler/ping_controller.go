package handler

import (
	"music-service/pkg/utils"
	"net/http"
)

// HandlePing
// @Summary Ping server
// @Description Send request to server
// @Tags ping
// @Accept  json
// @Produce  json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (h *Handler) HandlePing(writer http.ResponseWriter, request *http.Request) {
	h.log.Info("from ping...")
	_ = utils.WriteJSON(writer, http.StatusOK, "pong")
}
