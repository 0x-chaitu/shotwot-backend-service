package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"shotwot_backend/internal/domain"
	"shotwot_backend/internal/service"
	jwtauth "shotwot_backend/pkg/auth"
	"shotwot_backend/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *Handler) initUsersRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/signup", h.userSignUp)
	r.Post("/signin", h.userSignIn)
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseUser)
		r.Put("/update", h.userUpdate)
		r.Delete("/delete", h.deleteUser)
		r.Mount("/briefapplications", h.initBriefApplicationsRoutes())
		r.Get("/{userId}", h.getUser)
	})
	return r

}

func (h *Handler) userSignUp(w http.ResponseWriter, r *http.Request) {
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
	tokens, err := h.services.Users.SignUp(r.Context(), inp)
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
	render.Status(r, http.StatusOK)
	render.Render(w, r, &TokenResponse{
		Tokens: tokens,
	})
}

func (h *Handler) userSignIn(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.services.Users.SignIn(r.Context(), inp)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Data:           res,
		Success:        true,
	})
}

func (h *Handler) userUpdate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user domain.User
	err := decoder.Decode(&user)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	user.Id = userIdentity.Subject
	updatedUser, err := h.services.Users.Update(ctx, &user)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, updatedUser)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	user, err := h.services.Users.GetUser(r.Context(), id)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, user)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIdentity := ctx.Value(userCtx{}).(*jwtauth.CustomClaims)
	id := userIdentity.Subject
	err := h.services.Users.Delete(ctx, id)
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
