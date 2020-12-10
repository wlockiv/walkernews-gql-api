package auth

import (
	"context"
	"github.com/wlockiv/walkernews/internal/tables"
	"github.com/wlockiv/walkernews/pkg/jwt"
	"net/http"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow authenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//	Validate jwt token
			tokenStr := header
			userId, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			//	create user and check if user exists
			users := tables.GetUserTable()
			result, err := users.GetById(userId)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user := &tables.User{
				ID:       userId,
				Username: result.Username,
			}

			// Add user to the context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// Call next with context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
