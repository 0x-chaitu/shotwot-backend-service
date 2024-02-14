package v1

import (
	"encoding/json"
	"net/http"
	"shotwot_backend/internal/domain"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initBriefsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		// r.Use(h.parseAdmin)
		r.Post("/create", h.createBrief)
		r.Get("/details/{briefId}", h.getBrief)
		r.Put("/update", h.briefUpdate)
		r.Post("/list", h.listBriefs)
		r.Delete("/{briefId}", h.deleteBrief)

	})
	return r

}

func (h *Handler) createBrief(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	if !(adminIdentity.AdminRole == jwtauth.SuperAdmin ||
		adminIdentity.AdminRole == jwtauth.Admin || adminIdentity.AdminRole == jwtauth.BriefManager) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "not enough permissions",
		})
	}
	decoder := json.NewDecoder(r.Body)
	var brief domain.BriefInput
	err := decoder.Decode(&brief)
	brief.Created = time.Now()
	brief.CreatedBy = adminIdentity.Subject
	if err != nil {
		logger.Errorf("Error in decoding brief %v", err)
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	resBrief, err := h.services.Briefs.Create(ctx, &brief)
	if err != nil {
		logger.Error(err)
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           resBrief,
	})
}

func (h *Handler) briefUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	if !(adminIdentity.AdminRole == jwtauth.SuperAdmin ||
		adminIdentity.AdminRole == jwtauth.Admin ||
		adminIdentity.AdminRole == jwtauth.BriefManager) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "not admin",
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	var brief domain.BriefInput
	err := decoder.Decode(&brief)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	updatedBrief, err := h.services.Briefs.Update(ctx, &brief)
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
		Data:           updatedBrief,
	})
}

func (h *Handler) listBriefs(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var predicate helper.BriefPredicate
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
		Success:        true,
		Data:           briefList,
	})
}

func (h *Handler) deleteBrief(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "briefId")
	if id == "" {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      domain.ErrInvalidInput.Error(),
		})
		return
	}
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	if !(adminIdentity.AdminRole == jwtauth.SuperAdmin ||
		adminIdentity.AdminRole == jwtauth.Admin) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "not admin",
		})
		return
	}

	err := h.services.Briefs.DeleteBrief(ctx, id)
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
	})
}

func (h *Handler) getBrief(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "briefId")

	if id == "" {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      domain.ErrInvalidInput.Error(),
		})
		return
	}

	ctx := r.Context()

	brief, err := h.services.Briefs.GetBrief(ctx, id)
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
		Data:           brief,
	})
}
