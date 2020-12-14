package auth

import (
	"context"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/pkg/jwt"
	"net/http"
	"os"
)

type contextKey struct {
	name string
}

type UserCtx struct {
	User    *model.User
	UserKey string
}

var userCtxKey = &contextKey{"user"}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//	Validate jwt token
			tokenStr := header
			claims, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				//next.ServeHTTP(w, r)
				return
			}

			userCtx := UserCtx{
				User: &model.User{
					ID:    claims["userId"],
					Email: claims["email"],
				},
				UserKey: claims["userKey"],
			}

			// Add user to the context
			ctx := context.WithValue(r.Context(), userCtxKey, &userCtx)

			// Call next with context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) (*UserCtx, error) {
	userCtx, ok := ctx.Value(userCtxKey).(*UserCtx)
	if !ok {
		// Return the basic client if no token
		clientCtx := UserCtx{
			User:    nil,
			UserKey: os.Getenv("FDB_SERVER_CLIENT_KEY"),
		}

		return &clientCtx, nil
	}
	return userCtx, nil
}
