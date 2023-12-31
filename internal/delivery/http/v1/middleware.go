package v1

import (
	"context"
	"net/http"
	"shotwot_backend/internal/domain"
	"strings"

	"github.com/go-chi/render"
)

const (
	auth = "Authorization"
)

type userCtx struct{}

func (h *Handler) parseUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(auth)
		splitToken := strings.Split(header, " ")
		if len(splitToken) != 2 {
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusForbidden,
				ErrorText:      domain.ErrNotAuthorized.Error(),
			})
			return
		}
		userId, err := h.services.Auth.UserIdentity(splitToken[1])
		if err != nil {
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
