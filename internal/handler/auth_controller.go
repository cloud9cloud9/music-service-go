package handler

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"music-service/internal/models"
	"music-service/internal/security"
	"music-service/pkg/utils"
	"net/http"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
	internalServerError   = errors.New("internal server error")
	errUnauthorized       = errors.New("user is unauthorized")
)

// HandleLogin
// @Summary User Login
// @Tags auth
// @Description User authentication
// @ID login
// @Accept  json
// @Produce  json
// @Param input body models.LoginDto true "User credentials"
// @Success 200 {object} map[string]interface{} "token"
// @Failure 400 {object} error "invalid parsing JSON"
// @Failure 500 {object} error "internal server error"
// @Router /api/v1/login [post]
func (h *Handler) HandleLogin(writer http.ResponseWriter, request *http.Request) {
	var input models.LoginDto

	if err := utils.ParseJSON(request, &input); err != nil {
		h.log.Error("HANDLER: error parsing JSON: ", err)
		utils.InvalidParsingJSON(writer)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		h.log.Error("HANDLER: error validating input: ", err)
		er := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid input: %v", er))
		return
	}

	user, err := h.services.Authorization.GetUserByEmail(input.Email)
	if err != nil {
		h.log.Error("HANDLER: error getting user by email: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, internalServerError)
		return
	}

	if !security.CompareHashAndPassword(user.Password, []byte(input.Password)) {
		h.log.Error("HANDLER: error comparing passwords: ", err)
		utils.WriteError(writer, http.StatusBadRequest, errInvalidCredentials)
		return
	}

	token, err := h.services.Authorization.CreateToken(user.Username, user.Password)
	if err != nil {
		h.log.Error("HANDLER: error creating token: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, internalServerError)
		return
	}

	h.log.Info("HANDLER: token created: ", token)
	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

// HandleRegister
// @Summary Register
// @Tags auth
// @Description User registration
// @ID register
// @Accept  json
// @Produce  json
// @Param input body models.RegisterDto true "Register credentials"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} error "invalid parsing JSON"
// @Failure 500 {object} error "internal Server Error"
// @Router /api/v1/register [post]
func (h *Handler) HandleRegister(writer http.ResponseWriter, request *http.Request) {
	var input models.RegisterDto

	if err := utils.ParseJSON(request, &input); err != nil {
		h.log.Error("HANDLER: error parsing JSON: ", err)
		utils.InvalidParsingJSON(writer)
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		h.log.Error("HANDLER: error validating input: ", err)
		errors := err.(validator.ValidationErrors)
		utils.WriteError(writer, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	_, err := h.services.Authorization.GetUserByEmail(input.Email)
	if err == nil {
		h.log.Error("HANDLER: error getting user by email: ", err)
		utils.WriteError(writer, http.StatusBadRequest,
			fmt.Errorf("User with email %s already exists", input.Email))
		return
	}

	hash, err := security.HashPassword(input.Password)
	if err != nil {
		h.log.Error("HANDLER: error hashing password: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, internalServerError)
		return
	}

	err = h.services.Authorization.CreateUser(&models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hash,
	})
	if err != nil {
		h.log.Error("HANDLER: error creating user: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, internalServerError)
		return
	}

	h.log.Info("HANDLER: user created: ", input)
	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}

// LogoutHandler
// @Summary Logout
// @Tags auth
// @Description User logout
// @ID logout
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string "status"
// @Failure 500 {object} error "internal Server Error"
// @Router /api/v1/logout [post]
func (h *Handler) LogoutHandler(writer http.ResponseWriter, request *http.Request) {
	userId, err := getUserId(request.Context())
	if userId == 0 || err != nil {
		h.log.Error("HANDLER: error getting user id: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, errUnauthorized)
		return
	}

	err = h.services.Authorization.InvalidateToken(userId)
	if err != nil {
		h.log.Error("HANDLER: error invalidating token: ", err)
		utils.WriteError(writer, http.StatusInternalServerError, internalServerError)
		return
	}

	utils.WriteJSON(writer, http.StatusOK, map[string]interface{}{
		"status": "successfully logged out",
		"id":     userId,
	})
}
