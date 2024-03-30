package handler

import (
	"context"
	"net/http"

	"github.com/OliviaDilan/async_arc/task/internal/auth"
)

func AuthMiddleware(authClient auth.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, "Missing Authorization", http.StatusUnauthorized)
				return
			}

			user, err := authClient.DecodeToken(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := WithUserContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

type ctxUserKey struct {}

func WithUserContext(ctx context.Context, user *auth.User) context.Context {
	return context.WithValue(ctx, ctxUserKey{}, user)
}

func UserFromContext(ctx context.Context) *auth.User {
	user, ok := ctx.Value(ctxUserKey{}).(*auth.User)
	if !ok {
		return nil
	}
	return user
}