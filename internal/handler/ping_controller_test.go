package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_HandlePing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful ping request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `"pong"`,
		},
		{
			name:           "incorrect HTTP method",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectedBody:   `"pong"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, ping, nil)

			http.HandlerFunc(handler.HandlePing).ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
