package auth

import (
	"context"
	"github.com/wlockiv/walkernews/graph/model"
	"github.com/wlockiv/walkernews/pkg/jwt"
	"net/http"
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

			// Allow authenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			//	Validate jwt token
			tokenStr := header
			claims, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			//	create user and check if user exists
			//users, err := controllers.GetUserTable()
			//if err != nil {
			//	next.ServeHTTP(w, r)
			//	return
			//}

			//result, err := users.GetById(userId)
			//if err != nil {
			//	next.ServeHTTP(w, r)
			//	return
			//}

			//user := &controllers.User{
			//	ID:       userId,
			//	Username: result.Username,
			//}

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

func ForContext(ctx context.Context) *UserCtx {
	raw := ctx.Value(userCtxKey).(*UserCtx)
	return raw
}
