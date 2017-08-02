package axe

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/rs/cors"
)

// CSRF set csrf key
func CSRF(key string, secure bool, hnd http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(key),
		csrf.RequestHeader("Authenticity-Token"),
		csrf.FieldName("authenticity_token"),
		csrf.CookieName("csrf"),
		csrf.Secure(secure),
	)(hnd)
}

// CORS cors header
func CORS(debug bool, hnd http.Handler, hosts ...string) http.Handler {
	return cors.New(cors.Options{
		AllowedHeaders:   []string{"Authorization", "Cache-Control", "X-Requested-With"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowCredentials: true,
		Debug:            debug,
		AllowedOrigins:   hosts,
	}).Handler(hnd)
}
