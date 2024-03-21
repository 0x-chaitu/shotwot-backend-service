package v1

import (
	"encoding/json"
	"net/http"
	"shotwot_backend/internal/domain"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/logger"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initAdminAssetsRoutes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseAdmin)
		r.Post("/create", h.adminCreateAsset)
		r.Post("/update", h.adminUpdateAsset)
		r.Post("/list", h.adminGetAssetList)
	})
	return r

}

func (h *Handler) adminCreateAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	userID := userIdentity.Subject

	decoder := json.NewDecoder(r.Body)
	var asset domain.AssetInput
	err := decoder.Decode(&asset)
	if asset.Asset == nil {
		asset.Asset = &domain.Asset{}
	}
	asset.Asset.Created = time.Now()
	asset.Asset.UserId = userID
	if err != nil {
		logger.Errorf("Error in decoding brief %v", err)
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	resAsset, err := h.services.Assets.Create(ctx, &asset)
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
		Data:           resAsset,
	})
}

func (h *Handler) adminUpdateAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var asset domain.Asset
	err := decoder.Decode(&asset)
	if err != nil {
		logger.Errorf("Error in decoding brief %v", err)
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	resAsset, err := h.services.Assets.Update(ctx, &asset)
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
		Data:           resAsset,
	})
}

func (h *Handler) adminGetAssetList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	var asset domain.Asset
	err := decoder.Decode(&asset)
	if err != nil {
		logger.Errorf("Error in decoding brief %v", err)
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	resAsset, err := h.services.Assets.GetAllAssets(ctx)
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
		Data:           resAsset,
	})
}
