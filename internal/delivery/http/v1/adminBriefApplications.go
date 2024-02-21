package v1

import (
	"encoding/json"
	"net/http"
	"shotwot_backend/internal/domain"
	jwtauth "shotwot_backend/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initBriefApplicationRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseAdmin)
		r.Get("/applications/{briefId}", h.listBriefApplications)
		r.Get("/application/{applicationId}", h.getBriefApplication)
		r.Put("/application/action", h.updateBriefApplication)

	})
	return r
}

func (h *Handler) listBriefApplications(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "briefId")
	briefList, err := h.services.BriefApplications.GetBriefApplications(r.Context(), id)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           briefList,
	})
}

func (h *Handler) getBriefApplication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "applicationId")
	application, err := h.services.BriefApplications.GetBriefApplication(r.Context(), id)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           application,
	})
}

func (h *Handler) updateBriefApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	if !(adminIdentity.AdminRole == jwtauth.SuperAdmin ||
		adminIdentity.AdminRole == jwtauth.Admin ||
		adminIdentity.AdminRole == jwtauth.Curator) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "not enough permissions",
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	var briefApplication *domain.BriefApplication
	err := decoder.Decode(&briefApplication)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	application, err := h.services.BriefApplications.UpdateBriefApplication(r.Context(), briefApplication)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           application,
	})
}
