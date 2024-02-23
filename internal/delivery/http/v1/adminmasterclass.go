package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shotwot_backend/internal/domain"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initAdminMasterClassRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseAdmin)
		r.Post("/playlist/create", h.createPlaylist)
		r.Post("/playlist/list", h.listPlaylist)
	})
	return r

}

// Admin : Create Playlist Api
func (h *Handler) createPlaylist(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Hekko")
	ctx := r.Context()

	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	if !(adminIdentity.AdminRole == jwtauth.SuperAdmin ||
		adminIdentity.AdminRole == jwtauth.Admin || adminIdentity.AdminRole == jwtauth.BriefManager) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "not enough permissions",
		})
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var playlist domain.PlaylistInput
	if err := decoder.Decode(&playlist); err != nil {
		logger.Errorf("Error in decoding playlist: %v", err)
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	playlist.Created = time.Now()
	playlist.CreatedBy = adminIdentity.Subject

	resPlaylist, err := h.services.MasterClass.CreatePlaylist(ctx, &playlist)
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
		Data:           resPlaylist,
	})
}

// Admin : Get Playlist Api
func (h *Handler) listPlaylist(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var predicate helper.PlaylistPredicate
	err := decoder.Decode(&predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	playlist, err := h.services.MasterClass.GetPlaylists(r.Context(), &predicate)
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
		Data:           playlist,
	})
}
