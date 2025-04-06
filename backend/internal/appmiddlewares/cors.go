package appmiddlewares

import (
	"net/http"

	"github.com/go-chi/cors"
)

type AppCorsMiddleware struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	allowCredentials bool
	maxAge           int
	debug            bool
}

func NewAppCorsMiddleware(
	allowedOrigins []string,
	allowedMethods []string,
	allowedHeaders []string,
	allowCrecentials bool,
	maxAge int,
	debug bool,
) *AppCorsMiddleware {
	return &AppCorsMiddleware{
		allowedOrigins:   allowedOrigins,
		allowedMethods:   allowedMethods,
		allowedHeaders:   allowedHeaders,
		allowCredentials: allowCrecentials,
		maxAge:           maxAge,
		debug:            debug,
	}
}

func (m *AppCorsMiddleware) Handler() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   m.allowedOrigins,
		AllowedMethods:   m.allowedMethods,
		AllowedHeaders:   m.allowedHeaders,
		AllowCredentials: m.allowCredentials,
		MaxAge:           m.maxAge,
		Debug:            m.debug,
	})
}
