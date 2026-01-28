package login

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"AuthServis/internal/domain/models"
	"AuthServis/internal/lib/api/response"
	"AuthServis/internal/lib/jwt"
	"AuthServis/internal/storage"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
	Token string `json:"token,omitempty"`
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

func New(log *slog.Logger, userProvider UserProvider, tokenTTL time.Duration, appSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.login.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.Any("error", validateErr))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		user, err := userProvider.User(r.Context(), req.Email)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Warn("user not found", slog.String("email", req.Email))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid email or password"))
				return
			}
			log.Error("failed to get user", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(req.Password)); err != nil {
			log.Info("invalid password", slog.String("email", req.Email))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid email or password"))
			return
		}

		token, err := jwt.New(user, appSecret, tokenTTL)
		if err != nil {
			log.Error("failed to generate token", slog.Any("error", err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Token:    token,
		})
	}
}
