package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"strings"
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
		expectedStatus int
		expectedBody   string
		mockSetup      func()
		authHeader     string
		isJSON         bool
	}{
		{
			name:           "successful ping",
			expectedStatus: http.StatusOK,
			expectedBody:   `"pong"`,
			authHeader:     "Bearer validToken",
			mockSetup: func() {
				mockAuthService.EXPECT().IsTokenValid("validToken").Return(true, nil).Times(1)
				mockAuthService.EXPECT().ParseToken("validToken").Return(1, nil).Times(1)
			},
			isJSON: true,
		},
		{
			name:           "invalid token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid token"}`,
			authHeader:     "Bearer invalidToken",
			mockSetup: func() {
				mockAuthService.EXPECT().IsTokenValid("invalidToken").Return(false, nil).Times(1)
			},
			isJSON: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, ping, nil)
			req.Header.Set(authHeader, tt.authHeader)

			tt.mockSetup()

			handler.userIdentity(http.HandlerFunc(handler.HandlePing)).ServeHTTP(rec, req)

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
