package v1

import (
	"encoding/json"
	"net/http"
	"shotwot_backend/internal/domain"
	jwtauth "shotwot_backend/pkg/auth"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// route is users/savedbriefs/create
func (h *Handler) initSavedBriefsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseUser)
		r.Post("/create", h.saveBrief)
		r.Post("/list", h.listSavedBriefs)
	})
	return r
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
	// userId := "1hyvi8oFC0SKUXjBFcf0Q9uH9Be2"

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
