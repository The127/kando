package middlewares

import (
	"kando-backend/log"
	"net/http"
)

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Logger.Infof("[%s] %v", r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}
