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

const (
	test = "/test"
)

func TestUserIdentity(t *testing.T) {
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
		authHeader     string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:           "Empty Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization Header Format",
			authHeader:     "Basic invalidToken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty Bearer Token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid Token",
			authHeader: "Bearer invalidToken",
			mockSetup: func() {
				mockAuthService.EXPECT().IsTokenValid("invalidToken").Return(false, nil)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:       "Success",
			authHeader: "Bearer validToken",
			mockSetup: func() {
				mockAuthService.EXPECT().IsTokenValid("validToken").Return(true, nil)
				mockAuthService.EXPECT().ParseToken("validToken").Return(1, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test, nil)
			if tt.authHeader != "" {
				req.Header.Set(authHeader, tt.authHeader)
			}

			rec := httptest.NewRecorder()

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, 1, r.Context().Value(userCtx))
				w.WriteHeader(http.StatusOK)
			})

			handler.userIdentity(nextHandler).ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestLogRequest(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := &Handler{
		log: logging.NewLogger(),
	}
	handler.logRequest(nextHandler).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, test, nil))
}
