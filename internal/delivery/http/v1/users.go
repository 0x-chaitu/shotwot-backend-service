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
	otplessAuthSdk "github.com/otpless-tech/otpless-auth-sdk"
)

const (
	clientID     = "ON23N0MCOV49Y8C5N1FLNTNLB7FB6RY8"
	clientSecret = "2eqddf2458gxo1q4tda1enn5plmqfe7z"
)

func (h *Handler) initUsersRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/signup", h.userSignUp)
	r.Post("/signin", h.userSignIn)
	r.Post("/otp", h.userOtp)
	r.Post("/verify-otp", h.verifyOtp)
	r.Route("/", func(r chi.Router) {
		r.Use(h.parseUser)
		r.Put("/update", h.userUpdate)
		r.Delete("/delete", h.deleteUser)
		r.Get("/{userId}", h.getUser)

		// Mount
		// r.Mount("/briefapplications", h.initBriefApplicationsRoutes())
		// r.Mount("/savedbriefs", h.initSavedBriefsRoutes())

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

func (h *Handler) userOtp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var otp domain.Otp
	err := decoder.Decode(&otp)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	logger.Info(otp)
	var result *otplessAuthSdk.SendOTPResponse
	if otp.Email == "" {
		req := otplessAuthSdk.SendOTPRequest{
			PhoneNumber: otp.Phone,
			Channel:     "SMS",
			Expiry:      120,
			OtpLength:   6,
		}
		result, err = otplessAuthSdk.SendOTP(req, clientID, clientSecret)
		if err != nil {
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}
	} else {
		req := otplessAuthSdk.SendOTPRequest{
			Email:     otp.Email,
			Channel:   "EMAIL",
			Expiry:    120,
			OtpLength: 6,
		}
		result, err = otplessAuthSdk.SendOTP(req, clientID, clientSecret)
		if err != nil {
			render.Render(w, r, &ErrResponse{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorText:      err.Error(),
			})
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           result,
	})
}

func (h *Handler) verifyOtp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var otp domain.Otp
	ctx := r.Context()
	err := decoder.Decode(&otp)
	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}

	_, err = otplessAuthSdk.VerifyOTP(otp.OrderId, otp.Otp, otp.Email, otp.Phone, clientID, clientSecret)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorText:      err.Error(),
		})
		return
	}
	var res *service.AuthResponse
	if otp.Email == "" {
		user := &domain.User{
			Mobile:   otp.Phone,
			UserName: otp.Phone,
		}
		res, err = h.services.Users.GetOrCreateByPhone(ctx, user)

	} else {
		user := &domain.User{
			UserName: otp.Email,
			Email:    otp.Email,
		}
		res, err = h.services.Users.GetOrCreateByEmail(ctx, user)

	}

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
		Data:           res,
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
	user.UserId = userIdentity.Subject
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
	render.Render(w, r, &AppResponse{
		HTTPStatusCode: http.StatusOK,
		Success:        true,
		Data:           user,
	})
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
