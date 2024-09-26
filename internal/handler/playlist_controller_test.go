package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"music-service/internal/models"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_HandleCreatePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: services,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		input          models.CreatePlaylistDto
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful playlist creation",
			input: models.CreatePlaylistDto{
				Name: "test playlist",
			},
			mockSetup: func() {
				services.EXPECT().CreatePlaylist(&models.Playlist{
					Name: "test playlist",
				}).Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok","id":1}`,
		},
		{
			name: "error creating playlist",
			input: models.CreatePlaylistDto{
				Name: "test playlist",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
			userId:         1,
			mockSetup: func() {
				services.EXPECT().CreatePlaylist(gomock.Any()).Return(int64(0), errors.New("internal server error")).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, playlist, bytes.NewBuffer(reqBody))

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			tt.mockSetup()

			handler.HandleCreatePlaylist(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestHandler_HandleGetPlaylistById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: services,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful playlist get",
			userId: 1,
			mockSetup: func() {
				services.EXPECT().GetPlaylistById(1, 1).Return(&models.Playlist{
					ID:   1,
					Name: "test playlist",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"playlist":{"id":1,"name":"test playlist","user_id":0}}`,
		},
		{
			name:   "error getting playlist",
			userId: 1,
			mockSetup: func() {
				services.EXPECT().GetPlaylistById(1, 1).Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get(playlistById, handler.HandleGetPlaylistById)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/playlist/1", nil)

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			tt.mockSetup()

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}

func TestHandler_HandleGetAllPlaylists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: services,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful playlist get",
			userId: 1,
			mockSetup: func() {
				services.EXPECT().GetAllPlaylists(1).Return([]*models.Playlist{
					{
						ID:   1,
						Name: "test playlist",
					},
					{
						ID:   2,
						Name: "test playlist 2",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"test playlist","user_id":0}, {"id":2,"name":"test playlist 2","user_id":0}]`,
		},
		{
			name:   "error getting playlist",
			userId: 1,
			mockSetup: func() {
				services.EXPECT().GetAllPlaylists(1).Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get(playlist, handler.HandleGetAllPlaylists)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/playlist", nil)

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			tt.mockSetup()

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}

func TestHandler_HandleUpdatePlaylistById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: services,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		input          models.UpdatePlaylistDto
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful playlist update",
			input: models.UpdatePlaylistDto{
				Name: "test playlist",
			},
			userId: 1,
			mockSetup: func() {
				services.EXPECT().UpdatePlaylistById(1, &models.Playlist{
					Name: "test playlist",
					ID:   1,
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1, "name":"test playlist", "user_id":0}`,
		},

		{
			name: "error updating playlist",
			input: models.UpdatePlaylistDto{
				Name: "test playlist",
			},
			userId: 1,
			mockSetup: func() {
				services.EXPECT().UpdatePlaylistById(1, &models.Playlist{
					Name: "test playlist",
					ID:   1,
				}).Return(errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Put(playlistById, handler.HandleUpdatePlaylistById)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/playlist/1", nil)

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			b, err := json.Marshal(tt.input)
			require.NoError(t, err)
			req.Body = io.NopCloser(bytes.NewBuffer(b))

			tt.mockSetup()

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}

func TestHandler_HandleDeletePlaylistById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: services,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		userId         int
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful playlist delete",
			userId: 1,
			mockSetup: func() {
				services.EXPECT().DeletePlaylistById(1, 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"playlist deleted"}`,
		},
		{
			name:   "error deleting playlist",
			userId: 1,
			mockSetup: func() {
				services.EXPECT().DeletePlaylistById(1, 1).Return(errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Delete(playlistById, handler.HandleDeletePlaylistById)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/playlist/1", nil)

			ctx := context.WithValue(req.Context(), userCtx, tt.userId)
			req = req.WithContext(ctx)

			tt.mockSetup()

			r.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.JSONEq(t, tt.expectedBody, string(body))
		})
	}
}
