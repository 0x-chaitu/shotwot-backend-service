package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/service"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/helper"
	"shotwot_backend/pkg/logger"

	"github.com/gocarina/gocsv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initAdminRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/signin", h.adminSignIn)
	r.Get("/totalusers", h.getTotalUsers)
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseAdmin)
		r.Get("/details/{adminId}", h.getAdmin)
		r.Post("/users/download", h.downloadUsers)
		r.Post("/create", h.createAdmin)
		r.Post("/users/list", h.getAllUsers)
		r.Post("/users/search", h.searchUsers)
		r.Put("/update", h.adminUpdate)
		r.Get("/list", h.getAllAdmin)
		r.Delete("/delete", h.deleteAdmin)
		r.Mount("/brief", h.initBriefsRoutes())

	})
	return r

}

func (h *Handler) getAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "adminId")
	ctx := r.Context()
	logger.Info(id)
	if id == "me" {
		adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
		id = adminIdentity.Subject
	}
	admin, err := h.services.Admins.GetAdmin(r.Context(), id)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, admin)
}

func (h *Handler) downloadUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var predicate helper.UsersPredicate
	err := decoder.Decode(&predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	userArray, err := h.services.Users.Download(r.Context(), &predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Add("Content-Disposition", `attachment;Shotwot_Users.csv`)
	if err = gocsv.Marshal(userArray, w); err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorText:      err.Error(),
		})
		return
	}
}

func (h *Handler) createAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	decoder := json.NewDecoder(r.Body)
	var inp service.AccountAuthInput
	err := decoder.Decode(&inp)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	if adminIdentity.AdminRole > inp.Role {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      domain.ErrNotAuthorized.Error(),
		})
		return
	}
	err = h.services.Admins.CreateAdmin(r.Context(), inp)
	if err != nil {
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusConflict,
				ErrorText:      domain.ErrAccountAlreadyExists.Error(),
			})
			return
		} else if errors.Is(err, domain.ErrEmailPasswordInvalid) {
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusUnprocessableEntity,
				ErrorText:      domain.ErrEmailPasswordInvalid.Error(),
			})
			return
		}
		logger.Error("error during signup ", err)
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

func (h *Handler) adminSignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var inp service.AccountAuthInput
	err := decoder.Decode(&inp)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	tokens, err := h.services.Admins.SignIn(r.Context(), inp)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &TokenResponse{
		Tokens: tokens,
	})
}

func (h *Handler) adminUpdate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var admin domain.Admin
	err := decoder.Decode(&admin)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	admin.Id = adminIdentity.Subject
	updatedAdmin, err := h.services.Admins.Update(ctx, &admin)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, updatedAdmin)
}

func (h *Handler) getAllAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	role := adminIdentity.AdminRole
	if !(role == jwtauth.Admin || role == jwtauth.SuperAdmin) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "action not permitted",
		})
		return
	}
	adminList, err := h.services.Admins.GetAllAdmins(r.Context())
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
		Data:           adminList,
	})
}

func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	role := adminIdentity.AdminRole
	if !(role == jwtauth.Admin || role == jwtauth.SuperAdmin) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "action not permitted",
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	var predicate helper.UsersPredicate
	err := decoder.Decode(&predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	userList, err := h.services.Users.GetUsers(r.Context(), &predicate)
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
		Data:           userList,
	})
}

func (h *Handler) searchUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	role := adminIdentity.AdminRole
	if !(role == jwtauth.Admin || role == jwtauth.SuperAdmin) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "action not permitted",
		})
		return
	}
	decoder := json.NewDecoder(r.Body)
	var predicate helper.UsersPredicate
	err := decoder.Decode(&predicate)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	userList, err := h.services.Users.SearchUsers(r.Context(), &predicate)
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
		Data:           userList,
	})
}

func (h *Handler) getTotalUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.services.Users.TotalUsers(ctx)
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
		Data:           users,
	})
}

func (h *Handler) deleteAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	adminIdentity := ctx.Value(adminCtx{}).(*jwtauth.CustomAdminClaims)
	id := adminIdentity.Subject
	role := adminIdentity.AdminRole
	if !(role == jwtauth.Admin || role == jwtauth.SuperAdmin) {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      "action not permitted",
		})
		return
	}
	err := h.services.Admins.Delete(ctx, id)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
	})
}
