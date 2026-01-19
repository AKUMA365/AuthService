package register

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"AuthServis/internal/storage"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.register.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))

			responseError(w, "failed to decode request", http.StatusBadRequest)
			return
		}

		log.Info("request body decoded", slog.String("email", req.Email))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", slog.Any("error", validateErr))

			responseError(w, "invalid request fields", http.StatusBadRequest)
			return
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", slog.Any("error", err))

			responseError(w, "internal error", http.StatusInternalServerError)
			return
		}

		uid, err := userSaver.SaveUser(r.Context(), req.Email, passHash)

		if err != nil {
			if errors.Is(err, storage.ErrUserExists) {
				log.Warn("user already exists", slog.String("email", req.Email))
				responseError(w, "user already exists", http.StatusConflict)
				return
			}

			log.Error("failed to save user", slog.Any("error", err))
			responseError(w, "failed to save user", http.StatusInternalServerError)
			return
		}

		log.Info("user saved", slog.Int64("uid", uid))

		responseOK(w, uid)
	}
}

func responseError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "msg": msg})
}

func responseOK(w http.ResponseWriter, uid int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "OK",
		"uid":    uid,
	})
}
