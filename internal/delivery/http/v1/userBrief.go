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

func (h *Handler) initUserBriefRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseUser)
		r.Post("/apply", h.createBriefApplication)
		r.Post("/save", h.saveBrief)
		r.Get("/list/saved", h.listSavedBriefs)
		r.Get("/list", h.listBriefsUser)
		r.Get("/list/applied", h.getUserAppliedBriefs)
	})
	return r
}

func (h *Handler) listBriefsUser(w http.ResponseWriter, r *http.Request) {
	active := true
	predicate := helper.BriefPredicate{
		IsActive: &active,
		Order:    -1,
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

func (h *Handler) getUserAppliedBriefs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	userID := userIdentity.Subject
	briefList, err := h.services.BriefApplications.GetUserBriefApplications(r.Context(), userID)
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

func (h *Handler) createBriefApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	userID := userIdentity.Subject

	var briefInput domain.BriefApplication
	if err := json.NewDecoder(r.Body).Decode(&briefInput); err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	briefInput.Created = time.Now()
	briefInput.UserId = userID

	createdBrief, err := h.services.BriefApplications.Create(ctx, briefInput)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}

	// Render success response
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           createdBrief,
	})
}

func (h *Handler) saveBrief(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	userId := userIdentity.Subject

	var briefInput domain.SavedBriefInput
	if err := json.NewDecoder(r.Body).Decode(&briefInput); err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	briefInput.SavedBrief.Created = time.Now()
	briefInput.SavedBrief.UserId = userId

	savedBrief, err := h.services.SavedBriefs.CreateOrUpdate(ctx, &briefInput)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
	savedBrief.Status = ToggleStatus(savedBrief.Status)

	// Render success response
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           savedBrief,
	})
}

func ToggleStatus(status bool) bool {
	return !status
}

func (h *Handler) listSavedBriefs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	userId := userIdentity.Subject

	briefs, err := h.services.SavedBriefs.GetSavedBriefs(ctx, userId)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}

	// Render success response
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           briefs,
	})
}
