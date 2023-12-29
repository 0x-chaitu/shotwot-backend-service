package v1

import "net/http"

// Render for All Responses
func (ts *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Response is a wrapper response structure
type TokenResponse struct {
	Tokens interface{} `json:"tokens"`
}
