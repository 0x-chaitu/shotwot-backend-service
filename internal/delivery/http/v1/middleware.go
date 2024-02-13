package v1

import (
	"context"
	"net/http"
	"shotwot_backend/internal/domain"
	"shotwot_backend/pkg/logger"
	"strings"

	"github.com/go-chi/render"
)

const (
	auth = "Authorization"
)

type userCtx struct{}
type adminCtx struct{}

func (h *Handler) parseUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(auth)
		splitToken := strings.Split(header, " ")
		if len(splitToken) != 2 {
			logger.Error("user invalid token ")
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusForbidden,
				ErrorText:      domain.ErrNotAuthorized.Error(),
			})
			return
		}
		userId, err := h.services.Auth.UserIdentity(splitToken[1])
		if err != nil {
			logger.Errorf("error in authenticating user %v", err)
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusForbidden,
				ErrorText:      domain.ErrNotAuthorized.Error(),
			})
			return
		}
		ctx := context.WithValue(r.Context(), userCtx{}, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) parseAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(auth)
		splitToken := strings.Split(header, " ")
		if len(splitToken) != 2 {
			logger.Error("admin invalid token ")
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusForbidden,
				ErrorText:      domain.ErrNotAuthorized.Error(),
			})
			return
		}
		adminId, err := h.services.AdminAuth.AdminIdentity(splitToken[1])
		if err != nil {
			logger.Errorf("error in authenticating user %v", err)
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusForbidden,
				ErrorText:      domain.ErrNotAuthorized.Error(),
			})
			return
		}
		ctx := context.WithValue(r.Context(), adminCtx{}, adminId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
