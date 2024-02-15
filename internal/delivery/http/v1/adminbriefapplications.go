package v1

import (
	"encoding/json"
	"net/http"
	"shotwot_backend/pkg/helper"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initAdminBriefApplicationsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseAdmin)
		r.Post("/list", h.listBriefApplications)

	})
	return r

}

func (h *Handler) listBriefApplications(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var predicate helper.BriefApplicationsPredicate
	err := decoder.Decode(&predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	briefList, err := h.services.BriefApplications.GetBriefApplications(r.Context(), &predicate)
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
