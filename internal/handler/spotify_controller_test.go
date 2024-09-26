package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"music-service/internal/models"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_HandleGetTracksFromPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockSong(ctrl)
	handler := &Handler{
		services: &service.Service{
			Song: services,
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
	}{
		{
			name:       "successful get tracks from playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				services.EXPECT().GetAllSongsFromPlaylist(1, 1).Return([]*models.Song{
					{
						ID:    "1",
						Title: "test song",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"album":"", "album_cover":"", "artist":"", "duration":0, "external_url":"", "id":"1", "popularity":0, "preview_url":"", "release_date":"", "title":"test song"}]`,
		},
		{
			name:       "error getting tracks from playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				services.EXPECT().GetAllSongsFromPlaylist(1, 1).Return(nil, errors.New("test error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"test error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get(trackFromPlayList, handler.HandleGetTracksFromPlaylist)

			tt.mockSetup()
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/tracks/playlist/%d", tt.playlistId), nil)

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			require.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestHandler_HandleDeleteTrackFromPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockSong(ctrl)
	handler := &Handler{
		services: &service.Service{
			Song: services,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		playlistId     int
		trackId        string
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "successful delete track from playlist",
			playlistId: 1,
			trackId:    "1",
			userId:     1,
			mockSetup: func() {
				services.EXPECT().DeleteSongFromPlaylist(1, 1, "1").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Track removed from playlist"}`,
		},
		{
			name:       "error deleting track from playlist",
			playlistId: 1,
			trackId:    "1",
			userId:     1,
			mockSetup: func() {
				services.EXPECT().DeleteSongFromPlaylist(1, 1, "1").Return(errors.New("test error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"test error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Delete(insertAndDeleteTrack, handler.HandleDeleteTrackFromPlaylist)

			tt.mockSetup()
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/tracks/%s/playlist/%d",
				tt.trackId, tt.playlistId), nil)

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			require.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
