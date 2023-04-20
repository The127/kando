package utils

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func GetSessionId(r *http.Request) uuid.UUID {
	cookie, err := r.Cookie("SessionId")
	if err != nil {
		return uuid.Nil
	}

	sessionId, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.Nil
	}

	return sessionId
}

func CreateSessionCookie(sessionId uuid.UUID) *http.Cookie {
	return &http.Cookie{
		Name:     "SessionId",
		Value:    sessionId.String(),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

func DeleteSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "SessionId",
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}
