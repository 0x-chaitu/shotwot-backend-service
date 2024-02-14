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

// route is users/briefapplications/create
func (h *Handler) initBriefApplicationsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseUser)
		r.Post("/create", h.createBriefApplication)
	})
	return r
}

func (h *Handler) createBriefApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	userID := userIdentity.Subject

	var briefInput domain.BriefApplicationInput
	if err := json.NewDecoder(r.Body).Decode(&briefInput); err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	briefInput.BriefApplication.Created = time.Now()
	briefInput.BriefApplication.UserId = userID

	createdBrief, err := h.services.BriefApplications.Create(ctx, &briefInput)
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
