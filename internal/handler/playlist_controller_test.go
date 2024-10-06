package handler

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestHandler_HandleCreatePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	playlistService := mock_service.NewMockPlayList(ctrl)

	handler := &Handler{
		services: &service.Service{
			PlayList: playlistService,
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
		isJSON         bool
	}{
		{
			name: "successful playlist creation",
			input: models.CreatePlaylistDto{
				Name: "test playlist",
			},
			userId: 1,
			mockSetup: func() {
				playlistService.EXPECT().CreatePlaylist(&models.Playlist{
					Name:   "test playlist",
					UserId: 1,
				}).Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"ok","id":1}`,
			isJSON:         true,
		},
		{
			name: "error creating playlist",
			input: models.CreatePlaylistDto{
				Name: "test playlist",
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
			mockSetup: func() {
				playlistService.EXPECT().CreatePlaylist(gomock.Any()).Return(int64(0),
					errors.New("internal server error")).Times(1)
			},
			userId: 1,
			isJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, playlist, bytes.NewBuffer(reqBody))
			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			http.HandlerFunc(handler.HandleCreatePlaylist).ServeHTTP(rec, req)

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

func TestHandler_HandleGetPlaylistById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	playlistService := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: playlistService,
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
			name:       "successful playlist get",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				playlistService.EXPECT().GetPlaylistById(1, 1).Return(&models.Playlist{
					ID:   1,
					Name: "test playlist",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"playlist":{"id":1,"name":"test playlist","user_id":0}}`,
			isJSON:         true,
		},
		{
			name:       "error getting playlist",
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				playlistService.EXPECT().GetPlaylistById(1, 1).Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/playlist/%d", tt.playlistId), nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("playlistId", fmt.Sprintf("%d", tt.playlistId))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			http.HandlerFunc(handler.HandleGetPlaylistById).ServeHTTP(rec, req)

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

func TestHandler_HandleGetAllPlaylists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	playlistService := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: playlistService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		userId         int
		expectedStatus int
		expectedBody   string
		mockSetup      func()
		isJSON         bool
	}{
		{
			name:   "successful playlist get",
			userId: 1,
			mockSetup: func() {
				playlistService.EXPECT().GetAllPlaylists(1).Return([]*models.Playlist{
					{
						ID:     1,
						Name:   "test playlist",
						UserId: 1,
					},
					{
						ID:     2,
						Name:   "test playlist 2",
						UserId: 1,
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"test playlist","user_id":1}, {"id":2,"name":"test playlist 2","user_id":1}]`,
			isJSON:         true,
		},
		{
			name:   "error getting playlist",
			userId: 1,
			mockSetup: func() {
				playlistService.EXPECT().GetAllPlaylists(1).Return(nil, errors.New("internal server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
			isJSON:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, playlist, nil)

			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			http.HandlerFunc(handler.HandleGetAllPlaylists).ServeHTTP(rec, req)

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

func TestHandler_HandleUpdatePlaylistById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	playlistService := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: playlistService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		input          models.UpdatePlaylistDto
		playlistId     int
		userId         int
		expectedStatus int
		expectedBody   string
		mockSetup      func()
		isJSON         bool
	}{
		{
			name: "successful playlist update",
			input: models.UpdatePlaylistDto{
				Name: "test playlist",
			},
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				playlistService.EXPECT().UpdatePlaylistById(1, &models.Playlist{
					Name: "test playlist",
					ID:   1,
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1, "name":"test playlist", "user_id":0}`,
			isJSON:         true,
		},
		{
			name: "error updating playlist",
			input: models.UpdatePlaylistDto{
				Name: "test playlist",
			},
			playlistId: 1,
			userId:     1,
			mockSetup: func() {
				playlistService.EXPECT().UpdatePlaylistById(1, &models.Playlist{
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
			requestBody, _ := json.Marshal(tt.input)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/playlist/%d", tt.playlistId), strings.NewReader(string(requestBody)))
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("playlistId", fmt.Sprintf("%d", tt.playlistId))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			http.HandlerFunc(handler.HandleUpdatePlaylistById).ServeHTTP(rec, req)

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

func TestHandler_HandleDeletePlaylistById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	playlistService := mock_service.NewMockPlayList(ctrl)
	handler := &Handler{
		services: &service.Service{
			PlayList: playlistService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		playlistId     int
		userId         int
		expectedStatus int
		expectedBody   string
		mockSetup      func()
		isJSON         bool
	}{
		{
			name:           "successful playlist delete",
			playlistId:     1,
			userId:         1,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1}`,
			mockSetup: func() {
				playlistService.EXPECT().DeletePlaylistById(1, 1).Return(nil)
			},
			isJSON: true,
		},
		{
			name:           "internal server error",
			playlistId:     1,
			userId:         1,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error"}`,
			mockSetup: func() {
				playlistService.EXPECT().DeletePlaylistById(1, 1).Return(errors.New("internal server error"))
			},
			isJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/playlist/%d", tt.playlistId), nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("playlistId", fmt.Sprintf("%d", tt.playlistId))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			http.HandlerFunc(handler.HandleDeletePlaylistById).ServeHTTP(rec, req)

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
