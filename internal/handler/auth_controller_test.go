package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"music-service/internal/models"
	"music-service/internal/security"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_HandleRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthorization := mock_service.NewMockAuthorization(ctrl)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthorization,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		input          models.RegisterDto
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			input: models.RegisterDto{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockAuthorization.EXPECT().GetUserByEmail("test@example.com").Return(nil, errUserNotFound)
				mockAuthorization.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"status": "success"},
		},
		{
			name: "user already exists",
			input: models.RegisterDto{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockAuthorization.EXPECT().GetUserByEmail("existing@example.com").Return(&models.User{}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "User with email existing@example.com already exists"},
		},
		{
			name: "error hashing password",
			input: models.RegisterDto{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockAuthorization.EXPECT().GetUserByEmail("test@example.com").Return(nil, errUserNotFound)
				mockAuthorization.EXPECT().CreateUser(gomock.Any()).Return(internalServerError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "internal server error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", register, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.HandleRegister(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var responseBody map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &responseBody)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestHandler_HandleLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash, _ := security.HashPassword("password123")

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	handler := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		input          interface{}
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "Success",
			input: models.LoginDto{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				user := &models.User{
					Username: "testuser",
					Password: hash,
				}
				mockAuthService.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
				mockAuthService.EXPECT().CreateToken("testuser", user.Password).Return("token123", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"token": "token123"},
		},
		{
			name:           "Error parsing JSON",
			input:          "{invalid_json}",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "User not found",
			input: models.LoginDto{
				Email:    "notfound@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockAuthService.EXPECT().GetUserByEmail("notfound@example.com").Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "Invalid credentials",
			input: models.LoginDto{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				user := &models.User{
					Username: "testuser",
					Password: hash,
				}
				mockAuthService.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "Error creating token",
			input: models.LoginDto{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				user := &models.User{
					Username: "testuser",
					Password: hash,
				}
				mockAuthService.EXPECT().GetUserByEmail("test@example.com").Return(user, nil)
				mockAuthService.EXPECT().CreateToken("testuser", user.Password).Return("", errors.New("token creation error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, login, bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.HandleLogin(rec, req)

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != nil {
				var responseBody map[string]string
				json.NewDecoder(res.Body).Decode(&responseBody)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestHandler_LogoutHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	h := &Handler{
		services: &service.Service{
			Authorization: mockAuthService,
		},
		log: logging.NewLogger(),
	}

	tests := []struct {
		name           string
		userId         int64
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:   "Success",
			userId: 1,
			mockSetup: func() {
				mockAuthService.EXPECT().IsTokenValid("validToken").Return(true, nil).Times(1)
				mockAuthService.EXPECT().ParseToken("validToken").Return(1, nil).Times(1)
				mockAuthService.EXPECT().InvalidateToken(1).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"status": "successfully logged out"},
		},
		{
			name:           "Error getting user ID (User ID is 0)",
			userId:         0,
			mockSetup:      func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "user is unauthorized"},
		},
		{
			name:   "Error invalidating token",
			userId: 1,
			mockSetup: func() {
				mockAuthService.EXPECT().IsTokenValid("validToken").Return(true, nil).Times(1)
				mockAuthService.EXPECT().ParseToken("validToken").Return(1, nil).Times(1)
				mockAuthService.EXPECT().InvalidateToken(1).Return(errors.New("token invalidation error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, logout, nil)
			req.Header.Set(authHeader, "Bearer validToken")

			if tt.userId != 0 {
				ctx := context.WithValue(req.Context(), userCtx, tt.userId)
				req = req.WithContext(ctx)
			}

			tt.mockSetup()

			rec := httptest.NewRecorder()

			if tt.userId == 0 {
				http.HandlerFunc(h.LogoutHandler).ServeHTTP(rec, req)
			} else {
				h.userIdentity(http.HandlerFunc(h.LogoutHandler)).ServeHTTP(rec, req)
			}

			res := rec.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedBody != nil {
				var responseBody map[string]string
				err := json.NewDecoder(res.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Error decoding response body: %v", err)
				}
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}
