package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"AuthServis/internal/lib/api/response"
	"AuthServis/internal/lib/jwt"

	"github.com/go-chi/render"
)

type ctxKey string

const CtxKeyUID ctxKey = "uid"

func New(log *slog.Logger, appSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.auth.New"

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.Error("authorization header is required"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.Error("invalid authorization header format"))
				return
			}

			tokenStr := parts[1]

			uid, err := jwt.Parse(tokenStr, appSecret)
			if err != nil {
				log.Warn("invalid token", slog.String("op", op), slog.Any("error", err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.Error("invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), CtxKeyUID, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
