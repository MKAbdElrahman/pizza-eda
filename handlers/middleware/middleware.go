package middleware

import (
	"context"
	"net/http"
	"pizza/handlers/errorhandler"

	"github.com/alexedwards/scs/v2"
)

type UserStore interface {
	Exists(id int) (bool, error)
}

func RequireLoggingInFirst(sessionManager *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !userIdInSessionData(sessionManager, r) {
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
				return
			}
			w.Header().Add("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		})
	}
}

// the key will not be part of the session data until the user sucessfully logs in
func userIdInSessionData(sessionManager *scs.SessionManager, r *http.Request) bool {
	return sessionManager.Exists(r.Context(), "authenticatedUserID")
}

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

func UserAuthenticatedInContext(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func SetAuthenticatedUserInContext(sessionManager *scs.SessionManager, userStore UserStore, errorHandler *errorhandler.CentralErrorHandler) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := sessionManager.GetInt(r.Context(), "authenticatedUserID")
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}
			exists, err := userStore.Exists(id)
			if err != nil {
				errorHandler.HandleInternalServerError(w, r, err, "internal server error")
				return
			}
			if exists {
				ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, exists)
				ctx = context.WithValue(ctx, contextKey("userID"), id)
				r = r.WithContext(ctx)

			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(contextKey("userID")).(int)
	return userID, ok
}

func GetUserIDFromAuthenticatedContext(ctx context.Context) int {
	userID, _ := ctx.Value(contextKey("userID")).(int)
	return userID
}
