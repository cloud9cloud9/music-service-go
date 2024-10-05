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
	authService := mock_service.NewMockAuthorization(ctrl)
	handler := &Handler{
		services: &service.Service{
			Song:          songService,
			Authorization: authService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		playlistId     int
		authHeader     string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
		isJSON         bool
	}{
		{
			name:       "successful get tracks from playlist",
			playlistId: 1,
			authHeader: "Bearer validToken",
			mockSetup: func() {
				authService.EXPECT().IsTokenValid("validToken").Return(true, nil)
				authService.EXPECT().ParseToken("validToken").Return(1, nil)
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
			name:       "invalid token",
			playlistId: 1,
			authHeader: "Bearer invalidToken",
			mockSetup: func() {
				authService.EXPECT().IsTokenValid("invalidToken").Return(false, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid token"}`,
			isJSON:         false,
		},
		{
			name:       "error getting tracks from playlist",
			playlistId: 1,
			authHeader: "Bearer validToken",
			mockSetup: func() {
				authService.EXPECT().IsTokenValid("validToken").Return(true, nil)
				authService.EXPECT().ParseToken("validToken").Return(1, nil)
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
			req.Header.Set("Authorization", tt.authHeader)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			tt.mockSetup()

			handler.userIdentity(http.HandlerFunc(handler.HandleGetTracksFromPlaylist)).ServeHTTP(rec, req)

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
	authService := mock_service.NewMockAuthorization(ctrl)
	handler := &Handler{
		services: &service.Service{
			Song:          songService,
			Authorization: authService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		playlistId     int
		trackId        string
		authHeader     string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
		isJSON         bool
	}{
		{
			name:       "successful delete track from playlist",
			playlistId: 1,
			trackId:    "1",
			authHeader: "Bearer validToken",
			mockSetup: func() {
				authService.EXPECT().IsTokenValid("validToken").Return(true, nil)
				authService.EXPECT().ParseToken("validToken").Return(1, nil)
				songService.EXPECT().DeleteSongFromPlaylist(1, 1, "1").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"1"}`,
			isJSON:         true,
		},
		{
			name:       "error invalid token",
			playlistId: 1,
			trackId:    "1",
			authHeader: "Bearer invalidToken",
			mockSetup: func() {
				authService.EXPECT().IsTokenValid("invalidToken").Return(false, errors.New("error invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"error invalid token"}`,
			isJSON:         false,
		},
		{
			name:       "error deleting track from playlist",
			playlistId: 1,
			trackId:    "1",
			authHeader: "Bearer validToken",
			mockSetup: func() {
				authService.EXPECT().IsTokenValid("validToken").Return(true, nil)
				authService.EXPECT().ParseToken("validToken").Return(1, nil)
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
			req.Header.Set("Authorization", tt.authHeader)

			tt.mockSetup()

			handler.userIdentity(http.HandlerFunc(handler.HandleDeleteTrackFromPlaylist)).ServeHTTP(rec, req)

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
