package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/v3ronez/memi/types"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := &types.User{
			Name:  "Henrique",
			Email: "Henrique@gmail.com",
		}

		ctx := context.WithValue(r.Context(), types.UserCtxKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
