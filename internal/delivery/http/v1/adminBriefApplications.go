package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initAdminBriefApplicationsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseAdmin)
		r.Get("/{briefId}", h.listBriefApplications)
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
