package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"music-service/internal/models"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_HandleGetTracksFromPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	songService := mock_service.NewMockSong(ctrl)
	handler := &Handler{
		services: &service.Service{
			Song: songService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		playlistId     int
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
		isJSON         bool
	}{
		{
			name:       "successful get tracks from playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				songService.EXPECT().GetAllSongsFromPlaylist(1, 1).Return([]*models.Song{
					{
						ID:    "1",
						Title: "test song",
					},
				}, nil)
			},
			isJSON:         true,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"album":"", "album_cover":"", "artist":"", "duration":0, "external_url":"", "id":"1", "popularity":0, "preview_url":"", "release_date":"", "title":"test song"}]`,
		},
		{
			name:       "error getting tracks from playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				songService.EXPECT().GetAllSongsFromPlaylist(1, 1).Return(nil, errors.New("test error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"test error"}`,
			isJSON:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("playlistId", fmt.Sprintf("%d", tt.playlistId))
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/tracks/playlist/%d", tt.playlistId), nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			http.HandlerFunc(handler.HandleGetTracksFromPlaylist).ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.isJSON {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}

func TestHandler_HandleDeleteTrackFromPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	songService := mock_service.NewMockSong(ctrl)
	handler := &Handler{
		services: &service.Service{
			Song: songService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		playlistId     int
		userId         int
		trackId        string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
		isJSON         bool
	}{
		{
			name:       "successful delete track from playlist",
			playlistId: 1,
			trackId:    "1",
			mockSetup: func() {
				songService.EXPECT().DeleteSongFromPlaylist(1, 1, "1").Return(nil)
			},
			userId:         1,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"1"}`,
			isJSON:         true,
		},
		{
			name:       "error deleting track from playlist",
			playlistId: 1,
			userId:     1,
			trackId:    "1",
			mockSetup: func() {
				songService.EXPECT().DeleteSongFromPlaylist(1, 1, "1").Return(errTrackNotFound)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"track not found"}`,
			isJSON:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("trackId", tt.trackId)
			chiCtx.URLParams.Add("playlistId", fmt.Sprintf("%d", tt.playlistId))

			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/tracks/%s/playlist/%d",
				tt.trackId, tt.playlistId), nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}
			tt.mockSetup()

			http.HandlerFunc(handler.HandleDeleteTrackFromPlaylist).ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.isJSON {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}
