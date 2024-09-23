package handler

import (
	"context"
	"errors"
	"music-service/pkg/utils"
	"net/http"
	"strings"
	"time"
)

const (
	authHeader = "Authorization"
	userCtx    = "userId"
)

var (
	errEmptyAuthHeader   = errors.New("empty auth header")
	errInvalidToken      = errors.New("invalid token")
	errUsersTokenIsEmpty = errors.New("user not found")
	errUserNotFound      = errors.New("user not found")
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authHeader)
		if header == "" {
			utils.WriteError(w, http.StatusUnauthorized, errEmptyAuthHeader)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteError(w, http.StatusUnauthorized, errInvalidToken)
			return
		}

		if len(headerParts[1]) == 0 {
			utils.WriteError(w, http.StatusUnauthorized, errUsersTokenIsEmpty)
			return
		}

		userId, err := h.services.Authorization.ParseToken(headerParts[1])
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.log.Info(
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
		next.ServeHTTP(w, r)
	})
}

func getUserId(ctx context.Context) (int, error) {
	id, ok := ctx.Value(userCtx).(int)
	if !ok {
		return 0, errUserNotFound
	}
	return id, nil
}
