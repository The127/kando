package middlewares

import (
	"database/sql"
	"github.com/google/uuid"
	"kando-backend/httpErrors"
	"kando-backend/ioc"
	"kando-backend/services"
	"kando-backend/utils"
	"net/http"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scope := GetScope(r.Context())
		rcs := ioc.Get[*services.RequestContextService](scope)
		db := ioc.Get[*sql.DB](scope)

		sessionId := utils.GetSessionId(r)

		if sessionId == uuid.Nil {
			rcs.Error(httpErrors.Unauthorized())
			return
		}

		var sessionExists bool
		err := db.QueryRow(`select exists(select 1 from "public"."sessions" where "id" = $1);`,
			sessionId).
			Scan(&sessionExists)
		if err != nil {
			rcs.Error(httpErrors.Unauthorized())
			return
		}

		next.ServeHTTP(w, r)
	})
}
