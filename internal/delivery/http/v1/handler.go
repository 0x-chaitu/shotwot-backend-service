package v1

import (
	"net/http"
	"shotwot_backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	services *service.Services
}

// Render for All Responses
func (rd *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Response is a wrapper response structure
type Response struct {
	Data interface{} `json:"data"`
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status,omitempty"` // user-level status message
	AppCode    int64  `json:"code,omitempty"`   // application-specific error code
	ErrorText  string `json:"error,omitempty"`  // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() http.Handler {
	r := chi.NewRouter()
	r.Mount("/accounts", h.initAccountsRoutes())

	return r
}
