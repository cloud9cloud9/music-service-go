package handler

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"music-service/internal/models"
	"music-service/internal/service"
	mock_service "music-service/internal/service/mocks"
	"music-service/pkg/logging"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleRegister(t *testing.T) {
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
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
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
