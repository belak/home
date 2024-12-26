package internal

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/belak/home/internal/middleware"
	"github.com/belak/home/models"
)

var ErrAuthInvalidToken = errors.New("invalid token")

const AuthCookieName = "home-auth-token"

// AuthTokenFunc takes a token and returns either a UserInfo object or an error.
//
// If ErrAuthInvalidToken is returned, a 401 will be returned to the user,
// otherwise a 500.
type AuthTokenFunc func(ctx context.Context, token string) (*models.UserInfo, error)

func ExtractUser(ctx context.Context) *models.UserInfo {
	if user, ok := ctx.Value(middleware.CurrentUserContextKey).(*models.UserInfo); ok {
		return user
	}

	return nil
}

func AuthRequired(unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if u := ExtractUser(r.Context()); u == nil {
				unauthorizedHandler(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type extractTokenFunc func(*http.Request) (string, bool)

func authMiddleware(extractToken extractTokenFunc, tokenFunc AuthTokenFunc, unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if val, ok := extractToken(r); ok {
				user, err := tokenFunc(r.Context(), val)
				if errors.Is(err, ErrAuthInvalidToken) {
					unauthorizedHandler(w, r)
					return
				} else if err != nil {
					panic(err)
				}

				ctx := context.WithValue(r.Context(), middleware.CurrentUserContextKey, user)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractCookie(cookieName string) extractTokenFunc {
	return func(r *http.Request) (string, bool) {
		c, err := r.Cookie(cookieName)
		if err == http.ErrNoCookie {
			return "", false
		}
		return c.Value, true
	}
}

func extractHeader(r *http.Request) (string, bool) {
	vals := r.Header.Values("Authorization")
	if len(vals) < 1 {
		return "", false
	}

	for _, val := range vals {
		split := strings.SplitN(val, " ", 2)
		if len(vals) != 2 {
			continue
		}

		if split[0] != "Token" {
			continue
		}

		return split[1], true
	}

	return "", false
}

func AuthSetCookie(w http.ResponseWriter, t string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    t,
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
	})
}

func AuthClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func AuthCookieMiddleware(tokenFunc AuthTokenFunc, unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler {
	return authMiddleware(extractCookie(AuthCookieName), tokenFunc, unauthorizedHandler)
}

func AuthHeaderMiddleware(tokenFunc AuthTokenFunc, unauthorizedHandler http.HandlerFunc) func(http.Handler) http.Handler {
	return authMiddleware(extractHeader, tokenFunc, unauthorizedHandler)
}
