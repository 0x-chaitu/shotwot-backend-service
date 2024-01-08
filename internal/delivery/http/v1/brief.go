package v1

import (
	"encoding/json"
	"net/http"
	"shotwot_backend/internal/domain"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/helper"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initBriefsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {

		r.Use(h.parseAdmin)
		r.Post("/create", h.createBrief)
		r.Put("/update", h.userUpdate)
		r.Post("/list", h.listBriefs)
		r.Delete("/delete", h.deleteUser)
	})
	return r

}

func (h *Handler) createBrief(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	if !(adminIdentity.AdminRole == jwtauth.SuperAdmin ||
		adminIdentity.AdminRole == jwtauth.Admin) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "not admin",
		})
	}
	decoder := json.NewDecoder(r.Body)
	var brief domain.Brief
	err := decoder.Decode(&brief)
	brief.Created = time.Now()
	brief.CreatedBy = adminIdentity.Subject
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	err = h.services.Briefs.Create(r.Context(), &brief)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		SuccessText:    "brief created",
	})
}

func (h *Handler) listBriefs(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var predicate helper.Predicate
	err := decoder.Decode(&predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	briefList, err := h.services.Briefs.GetBriefs(r.Context(), &predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Data:           briefList,
	})
}
