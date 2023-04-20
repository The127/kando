package middlewares

import (
	"context"
	"github.com/gorilla/mux"
	"kando-backend/ioc"
	"net/http"
)

func ScopeMiddleware(dp *ioc.DependencyProvider) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scope := dp.NewScope()
			defer scope.Close()

			r = r.WithContext(context.WithValue(r.Context(), "scope", scope))
			next.ServeHTTP(w, r)
		})
	}
}

func GetScope(ctx context.Context) *ioc.DependencyProvider {
	return ctx.Value("scope").(*ioc.DependencyProvider)
}
